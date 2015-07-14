module PAC__TestCases_GetCommitMessageOnCorrectBranch
  require 'pp'
  require 'fileutils'
  require 'zip/zip'
  gem 'test-unit'
  require 'test/unit'
  require 'ci/reporter/rake/test_unit_loader.rb'
  require 'fileutils'
  require 'open3'

  # This test class have tests that verifies that we collect the correct commits
  # for a given repository called pac-test-repo.zip.
  # The test repository is contructed in a way that holds branches that are not merged
  # in yet, as well as merged branches. Merges are done like the Pretested Integration Plugin
  # for the Automated Git flow would do it (git merge -no-ff).
  # These verifications are do ensure the tree-walker traversing the git commits
  # follows the correct path in the topological search.
  #
  # These test cases are created based on a specific user case, where an issue was found
  # if the history had certain topoplog. Problem was the walker collecting the git commits
  # wasn't following topological order, but data order. Thus the pac-test-repo repository
  # created shows exactly this problem.
  #
  # Every test have three phases, each one verifies two lists - whitelist and blacklist
  # 1. Verify whitelist, that is all commits on the branch (master) between topological 
  # order of from and to SHA givens is collected. Blacklist is those commits not expected to be found, those only belonging to a branch that is not master
  # 2. Verifies whitelist of issue numbers (id, or references they are also called) which is a list of references found in the commit message. 
  # The blacklist are those issue numbers we do not expect because they are in commits that should be found.
  # 3. The un-referenced verification is checking that all commits that doesn't have a reference is in our list of such commits, while there should be no
  # commits there if they have a reference. The un-references list of commits is typically presented in the report.
  class TopologyWalkerIssueTest < Test::Unit::TestCase

    # Order in which to run the test.
    # Order matters here as we need to run simulator first to get logfile
    # parse
    self.test_order = :alphabetic

    # Helper function to unzip our repo
    def unzip_file (file, destination)
      Zip::ZipFile.open(file) { |zip_file|
        zip_file.each { |f|
          f_path=File.join(destination, f.name)
          FileUtils.mkdir_p(File.dirname(f_path))
          zip_file.extract(f, f_path) unless File.exist?(f_path)
        }
        puts "Successfully extracted"
      }
    end

    # Test setup: executed before each test
    def setup
      unzip_file("test/resources/pac-test-repo.zip", "test/resources")
    end

    # Executed after each test
    def teardown
      FileUtils.rm_rf("test/resources/pac-test-repo")
    end
    
    # Verifies that all commits on the master are found correct and only commit from master.
    def test_check_commit_found_on_master
      pp "******************************************************"
      pp " test_check_commit_found_on_master"
      pp "******************************************************"
      require_relative '../../../lib/core'
      settings_file = File.join(File.dirname(__FILE__), '../../resources/pac-test-repo_settings.yml')
      settings = YAML::load(File.open(settings_file))
      
      
      branch="master"
      # make sure we are on the expected branch as tests run abitrary orders
      system("cd test/resources/pac-test-repo && git checkout #{ branch}")


      Core.load(settings)


      # both these are on master, first and last commit on master
      # ./pac.rb --sha <to> []<from>]
      from="d926b6bf510abcda2ceff4ad01693a694e65141a"
      to="c1ece74c3411d0f19c49a1b193d3a8aec7376aa1"  ## note that c16bd71 is not last commit, last one is b437bdc - this affect what we verifies below :-)
      commit_map = Core.vcs.get_commit_messages_by_commit_sha(to,from)
           
      pp "########################################################################################"
      pp "All commit found between 'from' SHA #{ from } and 'to' SHA #{ to } on branch #{ branch }"
      pp "is in the commit map:"
      pp "########################################################################################"
      pp commit_map
      pp "########################################################################################"

      expected_SHAs = [
        "d926b6bf510abcda2ceff4ad01693a694e65141a",
        "f5bcb543a73e7d168f1c99c1da1ae0b9f84f3543",
        "3e2b00ad7b844d4708d7320788d2b54f917e4640",
        "1fbeb9eda1d620d44208c919f326f7e16485b16f",
        "e364d0e8c58efcc1069d516f83f14660df0a7dcd",
        "24cf861d3ca6da1af780c72cd39ca16f5b90dc56",
        "9439c822426bff6ac8bfebb66a893a281872ba38",
        "7088338c011cbcfaa69cb3a14bd12553f9b42584",
        "c1ece74c3411d0f19c49a1b193d3a8aec7376aa1",
        "4069c1010c7536acb584411fd71657eeb27eda65"
      ]

      # All branches are merged to master
      not_expected_SHAs = []

      pp "Checking with test asserts if the SHAs expected are found..."
      # based on our created test repository
      expected_SHAs.each do |sha|
        assert_true(commit_map.has_key?(sha), "Commit map didn't contain the expected SHA which is on #{ branch }: " + sha)
      end
      not_expected_SHAs.each do |sha|
        assert_false(commit_map.has_key?(sha), "Commit map included SHA that was NOT expected (does not exist on #{ branch }): " + sha)
      end
      pp "DONE - asserts okay for expected SHAs"
      



      grouped_by_task_id = Core.task_system.task_id_list(commit_map)
      pp "########################################################################################"
      pp "List of all task ids (references) in the commits just found:"
      pp "########################################################################################"
      pp grouped_by_task_id
      pp "########################################################################################"

      pp "Checking with test asserts if the id expected are there:"
      # expected ids are all those task ids that matches our regexp in the configuration file
      # but only in those commits in the commit map above
      expected_ids = [
        "1"
      ]
      # obvious these ids are from the commits not found, thus the these task ids should be found either
      not_expected_ids = [ ] # There are no un-expected ID as there is only one id and it is expect to eb found

      expected_ids.each do |id|
        assert_true(grouped_by_task_id.has_key?(id), "Didn't find task reference for '#{ id }' as expected")
      end
      not_expected_ids.each do |id|
        assert_false(grouped_by_task_id.has_key?(id), "Found task reference for '#{ id }' which was NOT expected")
      end
      pp "DONE - asserts okay for expected ids"




      unreferenced = Core.task_system.get_shas_without_reference(commit_map, grouped_by_task_id)
      pp "########################################################################################"
      pp "List of all commit that doesn't have a task ids (references):"
      pp "########################################################################################"
      pp unreferenced
      pp "########################################################################################"

      # These commits does not have task id references, but are found on the branch we traverse
      expected_unreferenced_SHAs = [
        "d926b6bf510abcda2ceff4ad01693a694e65141a",
        "f5bcb543a73e7d168f1c99c1da1ae0b9f84f3543",
        "3e2b00ad7b844d4708d7320788d2b54f917e4640",
        "1fbeb9eda1d620d44208c919f326f7e16485b16f",
        "e364d0e8c58efcc1069d516f83f14660df0a7dcd",
        "24cf861d3ca6da1af780c72cd39ca16f5b90dc56",
        "7088338c011cbcfaa69cb3a14bd12553f9b42584",
        "c1ece74c3411d0f19c49a1b193d3a8aec7376aa1",
        "4069c1010c7536acb584411fd71657eeb27eda65"
      ]
      
      # These are all commits on the branch we traverse, without a task id that matches our
      # regexp in the configuration file, as well as all those commits (with or without task ids
      # that is not on the branch
      referenced_SHAs = [
        "9439c822426bff6ac8bfebb66a893a281872ba38" 
      ]

      pp "Checking with test asserts for un-referenced commits :"
      # Checking on unreferenced shas:
      expected_unreferenced_SHAs.each do |sha|
        assert_true(unreferenced.include?(sha), "Found a SHA in unreferenced SHAs, those without issue reference, that was not expected: " + sha)
      end      
      referenced_SHAs.each do |sha|
        assert_false(unreferenced.include?(sha), "Found a SHA in unreferenced SHAs, those without issue reference, that was not expected: " + sha)
      end
      pp "DONE - assert okay for un-referenced commits"

    end

  end # class
end # module

