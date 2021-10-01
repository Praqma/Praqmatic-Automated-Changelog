# encoding: utf-8
require "rugged"
require_relative "model"
require_relative "logging"

module Vcs
  
  class GitVcs
    attr_accessor :repository
    attr_accessor :settings
    
    def initialize(settings)
      @settings = settings
      @repository = Rugged::Repository.new(@settings[:repo_location])    
    end
    
    def createWalker
      if @walker.nil?
        return Rugged::Walker.new(@repository)
      end
  
      @walker
    end
  
    def walker=(v)
      @walker = v
    end

    #Method that returns the newest commit a tag points to
    def get_latest_tag(treeish) 
      tag_collection = Rugged::TagCollection.new(repository)
      candidates = []
      tag_collection.each(treeish) do |tag|
        candidates << tag
      end

      if candidates.empty? 
        raise "[PAC] No matching tags found with approximation #{treeish}"
      end

      candidates.sort! {|a,b| a.target.time <=> b.target.time }

      Logging.verboseprint(0, "[PAC] Found latest tag: #{candidates.last.name}")

      candidates.last.name
    end

    #Super simplified query for git
    def get_delta(oldest, newest=nil)
      if newest.nil?
        head = repository.lookup(repository.head.target.oid)
      else
        head = repository.rev_parse(newest)
      end

      tail = repository.rev_parse(oldest)

      walker = createWalker

      commits = Model::PACCommitCollection.new
      #puts head.inspect
      walker.push(head.oid)
      walker.hide(tail.oid)

      Logging.verboseprint(0, "[PAC] Warning! Both :sparse and :filter_paths are defined! Using sparse") if @settings[:filter_paths] and @settings[:sparse]

      git_data = if @settings[:sparse]
                  `git log --pretty=format:"%h" --sparse`.split("\n")
                 elsif @settings[:filter_paths]
                   `git log --pretty=format:"%h" -- #{@settings[:filter_paths].join(' ')}`.split("\n")
                 else
                   nil
                 end

      walker.inject([]) do |c, commit|
        if not git_data or git_data.any? { |c|  commit.oid.to_s.include?(c) }
          p_commit = Model::PACCommit.new(commit.oid, commit.message, commit.time)
          Logging.verboseprint(3, "[PAC] Added commit #{commit.oid}")
          commits.add(p_commit)
        end
      end
      commits
    end
  end
  
end
