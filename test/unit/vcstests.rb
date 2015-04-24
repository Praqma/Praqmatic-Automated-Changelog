if( ENV['COVERAGE'] == 'on' )
  require 'simplecov'
  require 'simplecov-rcov'
  SimpleCov.formatter = SimpleCov::Formatter::RcovFormatter
  SimpleCov.start
end

require 'test/unit'
require 'flexmock/test_unit'
require 'rugged'
require_relative '../../lib/core'

class VcsTests < Test::Unit::TestCase
  def setup
    @mockCommits = {
    :commit_5 => flexmock('commit_5', :on, Rugged::Commit,
                            :oid => '555',
                            :message => 'Commit 5 message',
                            :time => Time.new(2013, 07, 03)
                           ),
    :commit_4 => flexmock('commit_4', :on, Rugged::Commit,
                            :oid => '444',
                            :message => 'Commit 4 message',
                            :time => Time.new(2013, 07, 01)
                           ),
    :commit_3 => flexmock('commit_3', :on, Rugged::Commit,
                            :oid => '333',
                            :message => 'Commit 3 message',
                            :time => Time.new(2013, 06, 22)
                           ),
    :commit_2 => flexmock('commit_2', :on, Rugged::Commit,
                            :oid => '222',
                            :message => 'Commit 2 message',
                            :time => Time.new(2013, 06, 15)
                           ),
    :commit_1 => flexmock('commit_1', :on, Rugged::Commit,
                            :oid => '111',
                            :message => 'Commit 1 message',
                            :time => Time.new(2013, 06, 05)
                           )
    }

    mockReferences = [
      flexmock('refs/remotes/origin/master', :on, Rugged::Reference,
               :name => 'refs/remotes/origin/master',
               :target => @mockCommits[:commit_5].oid
              ), 
      flexmock('refs/heads/master', :on, Rugged::Reference,
               :name => 'refs/heads/master',
               :target => @mockCommits[:commit_5].oid
              ), 
      flexmock('refs/tags/tag-1', :on, Rugged::Reference,
               :name => 'refs/tags/tag-1',
               :target => @mockCommits[:commit_2].oid
              ), 
      flexmock('refs/tags/tag-2', :on, Rugged::Reference,
               :name => 'refs/tags/tag-2',
               :target => @mockCommits[:commit_3].oid
              ), 
      flexmock('refs/tags/tag-3', :on, Rugged::Reference,
               :name => 'refs/tags/tag-3',
               :target => @mockCommits[:commit_4].oid
              ), 
      flexmock('refs/remotes/origin/HEAD', :on, Rugged::Reference,
               :name => 'refs/remotes/origin/HEAD'
              )
    ]

    mockTags = [
      flexmock('tag-1', :on, Rugged::Tag,
               :name => 'tag-1',
               :target => @mockCommits[:commit_2]
              ), 
      flexmock('tag-3', :on, Rugged::Tag,
               :name => 'tag-3',
               :target => @mockCommits[:commit_4]
              )
    ]

    @mockRepository = flexmock('repository', :on, Rugged::Repository,
                              :last_commit => @mockCommits[:commit_5]
                             )

    @mockRepository.should_receive(:head)
                  .and_return do |ref|
                    mockReferences.find { |ref| ref.name == 'refs/heads/master' }
                  end

    @mockRepository.should_receive(:refs)
                  .with(Regexp)
                  .and_return do |regexp|
                    mockReferences.select { |ref| ref.name.match(regexp) }
                  end

    @mockRepository.should_receive(:lookup)
                  .with(String)
                  .and_return do |oid|
                    commit = @mockCommits.values.find { |v| v.oid == oid }

                    if commit.nil?
                      raise Rugged::OdbError.new
                    end

                    commit
                  end

    @mockWalker = flexmock('walker', :on, Rugged::Walker)
    @mockWalker.should_receive(:each).with(Proc)
               .and_return do |block| 
                 @mockCommits.values.each { |c| block.call c } 
               end

    Vcs.repository = @mockRepository
    Vcs.walker = @mockWalker
  end

  def test_should_raise_ArgumentError_if_tail_tag_is_nil
    tailTag = nil

    assert_raise(ArgumentError) { Vcs.get_commit_messages_by_tag_name(tailTag) }
  end

  def test_should_raise_ArgumentError_if_tail_SHA_is_nil
    tailSHA = nil

    assert_raise(ArgumentError) { Vcs.get_commit_messages_by_commit_sha(tailSHA) }
  end

  def test_should_raise_ArgumentError_if_tail_time_is_nil
    tailTime = nil

    assert_raise(ArgumentError) { Vcs.get_commit_messages_by_commit_times(tailTime) }
  end

  def test_should_raise_ArgumentError_if_tail_tag_can_not_be_found
    tailTag = 'Nonexistant tail tag'

    assert_raise(ArgumentError) { Vcs.get_commit_messages_by_tag_name(tailTag) }
  end

  def test_should_raise_ArgumentError_if_tail_SHA_can_not_be_found
    tailSHA = 'Nonexistant tail SHA'

    assert_raise(ArgumentError) { Vcs.get_commit_messages_by_commit_sha(tailSHA) }
  end

  def test_should_raise_ArgumentError_if_head_tag_can_not_be_found
    headTag = 'Nonexistant head tag'
    tailTag = 'tag-1'

    assert_raise(ArgumentError) { Vcs.get_commit_messages_by_tag_name(tailTag, headTag) }
  end

  def test_should_raise_ArgumentError_if_head_SHA_can_not_be_found
    headSHA = 'Nonexistant head SHA'
    tailSHA = @mockCommits[:commit_3].oid

    assert_raise(ArgumentError) { Vcs.get_commit_messages_by_commit_sha(tailSHA, headSHA) }
  end

  def test_can_return_commit_messages_between_designated_tail_and_head_tags
    head = @mockCommits[:commit_4]
    headIndex = @mockCommits.values.find_index(head)
    availableElements = @mockCommits.values.slice(headIndex, @mockCommits.values.length)

    expectedCommitMessages = ["Commit 4 message", "Commit 3 message", "Commit 2 message"]
    headTag = 'tag-3'
    tailTag = 'tag-1'

    @mockWalker.should_receive(:push)
      .with(String)
      .once.and_return do
        @mockCommits = @mockCommits.select do |k, v|
          if availableElements.include? v
            { k => v }
          end
        end
      end

    commitMessages = Vcs.get_commit_messages_by_tag_name(tailTag, headTag)

    assert_equal(expectedCommitMessages, commitMessages)
  end

  def test_can_return_commit_messages_between_designated_tail_and_head_SHAs
    head = @mockCommits[:commit_4]
    headIndex = @mockCommits.values.find_index(head)
    availableElements = @mockCommits.values.slice(headIndex, @mockCommits.values.length)

    expectedCommitMessages = ["Commit 4 message", "Commit 3 message"]
    headSHA = @mockCommits[:commit_4].oid
    tailSHA = @mockCommits[:commit_3].oid

    @mockWalker.should_receive(:push)
      .with(String)
      .once.and_return do 
        @mockCommits = @mockCommits.select do |k, v|
          if availableElements.include? v
            { k => v }
          end
        end
      end

    commitMessages = Vcs.get_commit_messages_by_commit_sha(tailSHA, headSHA)

    assert_equal(expectedCommitMessages, commitMessages)
  end

  def test_can_return_commit_messages_between_designated_tail_and_head_times
    expectedCommitMessages = ["Commit 4 message", "Commit 3 message", "Commit 2 message"]
    headTime = @mockCommits[:commit_4].time
    tailTime = @mockCommits[:commit_2].time

    @mockWalker.should_receive(:push).with(String).once 

    commitMessages = Vcs.get_commit_messages_by_commit_times(tailTime, headTime)

    assert_equal(expectedCommitMessages, commitMessages)
  end

  def test_should_use_head_of_master_if_head_tag_param_is_nil
    expectedCommitMessages = ["Commit 5 message", "Commit 4 message", "Commit 3 message", "Commit 2 message"]
    headTag = nil
    tailTag = 'tag-1'

    commitMessages = Vcs.get_commit_messages_by_tag_name(tailTag, headTag)

    assert_equal(expectedCommitMessages, commitMessages)
  end

  def test_should_use_head_of_master_if_head_SHA_param_is_nil
    expectedCommitMessages = ["Commit 5 message", "Commit 4 message", "Commit 3 message"]
    headSHA = nil
    tailSHA = @mockCommits[:commit_3].oid

    commitMessages = Vcs.get_commit_messages_by_commit_sha(tailSHA, headSHA)

    assert_equal(expectedCommitMessages, commitMessages)
  end

  def test_should_use_head_of_master_if_head_time_parm_is_nil
    expectedCommitMessages = ["Commit 5 message", "Commit 4 message", "Commit 3 message"]
    headTime = nil
    tailTime = @mockCommits[:commit_3].time

    commitMessages = Vcs.get_commit_messages_by_commit_times(tailTime, headTime)

    assert_equal(expectedCommitMessages, commitMessages)
  end

#  def test_getCorrectTaskNumbers
#    changeset = ["Fixed Ticket#1", "fixed ticket#2", "This is a new commit messge with ticket#nOne","Commit without task reference"]
#    stats = Core::TracStatistics.new changeset
#    assert_equal(3, stats.taskReferenceCount)
#    assert_equal(1, stats.noneTaskCount)
#    assert_equal(4, stats.totalCommits)
#    assert_equal(1, stats.unspecifiedCommitCount)       
#  end
end
