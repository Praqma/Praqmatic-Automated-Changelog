#!/usr/bin/env ruby
# encoding: utf-8
require "mercurial-ruby"

module Vcs

  class MercurialVcs
    attr_accessor :repository
    def initialize(settings)
      Mercurial.configure do |conf|
        conf.hg_binary_path = "/usr/bin/hg"
      end
      @repository = Mercurial::Repository.open(settings['repo_location'])
    end

    #Converting from string to time and striping the hours
    def date_converter(datestring)
      DateTime.strptime(datestring, "%Y-%m-%d").to_time
    end

    #Converting from string to time object for the 'get_commits_by_time_with_hours' method
    def date_converter_complex(datestring)
      DateTime.strptime(datestring, "%Y-%m-%d %H:%M:%S").to_time
    end

    #Used in the 'get_commit_messages_by_tag_name' method
    def get_commits_by_time_with_hours(tailTime, headTime=nil)

      raise ArgumentError, 'Tail time parameter is nil' if tailTime.nil?

      headCommit = @repository.commits.tip
      headCommitDate = headCommit.date

      if headTime.nil?
      headTime = headCommitDate.to_time
      else
        headTime = date_converter_complex(headTime)
      end
      tailTime = date_converter_complex(tailTime)

      commit_messages = []

      @repository.commits.each do |commit|
        if headTime >= date_converter_complex(commit.date.to_s) && date_converter_complex(commit.date.to_s)  >= tailTime
          commit_messages.push(commit.message)
        end
      end
      commit_messages
    end

    #Get commits by specifying time span. Pattern example: ('2013-10-03', '2013-10-07')
    def get_commit_messages_by_commit_times (tailTime, headTime = nil)
      raise ArgumentError, 'Tail time parameter is nil' if tailTime.nil?

      headCommit = @repository.commits.tip
      headCommitDate = headCommit.date

      if headTime.nil?
        headTime = headCommitDate.to_time
      else
        headTime = date_converter(headTime.to_s)
      end
      tailTime = date_converter(tailTime.to_s)

      commit_messages = []

      @repository.commits.each do |commit|
        if headTime >= date_converter(commit.date.to_s) && date_converter(commit.date.to_s)  >= tailTime
          commit_messages.push(commit.message)
        end
      end

      commit_messages
    end

    #Get commit messages by specifying the secured hash algorithm (SHA) of the commits
    def get_commit_messages_by_commit_sha(tailCommitSHA, headCommitSHA=nil)
      raise ArgumentError, 'Tail commit SHA parameter is nil' if tailCommitSHA.nil?

      headCommit = @repository.commits.tip
      headCommitSha = headCommit.to_s

      if headCommitSHA.nil?
      headCommitSHA = headCommitSha
      end
      shas = []
      @repository.commits.each do |commit|

        if commit.hash_id.include? tailCommitSHA
          commitTailCommitSHADate = commit.date.to_s
          @commitTailCommitSHADate = commitTailCommitSHADate
        end
        if commit.hash_id.include? headCommitSHA
          commitHeadCommitSHADate = commit.date.to_s
          @commitHeadCommitSHADate = commitHeadCommitSHADate
        end
      end
      result = get_commits_by_time_with_hours(@commitTailCommitSHADate,@commitHeadCommitSHADate)
      return result
    end

    #Get commits by specifying a start and an end tag
    def get_commit_messages_by_tag_name(tailTagName, headTagName=nil)
      raise ArgumentError, 'Tail tag name is nil' if tailTagName.nil?

      commit_messages = []
      @repository.commits.each do |commit|
        if commit.tags_names.include? tailTagName
          commitTailTagDate = commit.date.to_s
          @commitTailTagDate = commitTailTagDate
        end
        if commit.tags_names.include? headTagName
          commitHeadTagDate = commit.date.to_s
          @commitHeadTagDate = commitHeadTagDate
        end
      end

      result = get_commits_by_time_with_hours(@commitTailTagDate,@commitHeadTagDate)
      return result
    end
  end

end
