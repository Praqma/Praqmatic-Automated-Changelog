# encoding: utf-8
require "rugged"
require_relative "model"

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
    
    def get_commit_messages_by_commit_times(tailTime, headTime=nil)
      raise ArgumentError, 'Tail time parameter is nil' if tailTime.nil?
  
      headTarget = @repository.head.target
      headCommit = @repository.lookup(headTarget)
  
      if headTime.nil?
        headTime = headCommit.time
      end
  
      headTime = headTime.to_i
      tailTime = tailTime.to_i
  
      walker = createWalker

      #We need to sort topologically (since commits can be ooo)
      walker.sorting(Rugged::SORT_TOPO)

      walker.push(headCommit.oid)
      
      commits = Model::PACCommitCollection.new

      walker.each do |commit|
        if headTime >= commit.time.to_i && commit.time.to_i >= tailTime
          p_commit = PACCommit.new(commit.oid, commit.message, commit.time)
          commits.add(p_commit)
        end
      end
       
      return commits
    end
  
    def get_commit_messages_by_commit_sha(tailCommitSHA, headCommitSHA=nil)
      raise ArgumentError, 'Tail commit SHA parameter is nil' if tailCommitSHA.nil?
  
      checkIfCommitExists = lambda do |sha|
        begin
          @repository.lookup(sha)
        rescue Rugged::OdbError
          raise ArgumentError, "Commit with SHA '#{sha}' can not be found" 
        end
      end
  
      checkIfCommitExists.(tailCommitSHA)
      
      if headCommitSHA.nil?
        headCommitSHA = @repository.head.target
      else
        checkIfCommitExists.(headCommitSHA)
      end
  
      walker = createWalker
      
      #We need to sort topologically (since commits can be ooo)
      walker.sorting(Rugged::SORT_TOPO)

      walker.push(headCommitSHA)
  
      commits = Model::PACCommitCollection.new
  
      begin
        walker.each do |commit|          
          if commit.oid == tailCommitSHA
            break
          end
          p_commit = Model::PACCommit.new(commit.oid, commit.message, commit.time)
          commits.add(p_commit)
        end
      rescue RuntimeError
        puts "Failed to load changeset"
      end
  
      return commits
    end
    
    def get_commit_messages_by_tag_name(tailTag, headTag=nil)
      raise ArgumentError, 'Tail tag parameter is nil' if tailTag.nil?
      
      #FIXME: How dow we do this properly? 
      getTaggedCommitTime = lambda do |tagName|
        if tailTag != "LATEST" 
          tagCommitTime = @repository.references["refs/tags/#{tagName}"]
              
          if tagCommitTime.nil?
            raise ArgumentError, "The tag with the name #{tailTag} not found in repository"
          end
          
          mycommit = @repository.lookup(tagCommitTime.target_id)

          #If target is another tag, get the time of that
          if (mycommit.class == Rugged::Tag::Annotation)
            mycommit.target.time
          else
            mycommit.time
          end
        else
          commitTimes = []
          release_regex = /tags/
          unless @settings[:release_regex].nil?
            release_regex =  Regexp.new @settings[:release_regex]  
          end

          tailTagReference = @repository.refs(release_regex).each do |ref|
            commitToSort = @repository.lookup(ref.target)                   
            if commitToSort.class == Rugged::Tag
              commitTimes << { :oid => ref.target, :time => commitToSort.target.time, :name  => commitToSort.name  }
            end 
          end
          
          raise ArgumentError, "No tags on current branch found" if commitTimes.empty?
          sorted = commitTimes.sort { |x,y| y[:time] <=> x[:time] }
          sorted.first[:time]        
        end
      end
  
      getTagReferenceOid = lambda do |tagName|
        if tailTag != "LATEST"
          tailTagReference = @repository.references["refs/tags/#{tagName}"]
          
          raise ArgumentError, "Tag reference '#{tagName}' can not be found" if tailTagReference.nil?
          tailTagReference
        else
          commitTimes = []
          release_regex = /tags/
          unless @settings[:release_regex].nil?
            release_regex =  Regexp.new @settings[:release_regex]  
          end
          tailTagReference = @repository.refs(release_regex).each do |ref|
            commitToSort = @repository.lookup(ref.target)
            
            if commitToSort.class == Rugged::Tag
              commitTimes << { :oid => ref.target, :time => commitToSort.target.time, :name  => commitToSort.name  }
            end
          end
          
          raise ArgumentError, "No tags on current branch found" if commitTimes.empty?
          sorted = commitTimes.sort { |x,y| y[:time] <=> x[:time] }
          sorted.first[:oid]
          
        end
      end
  
      getCommitOid = lambda do |tagOid|
        element = @repository.lookup(tagOid.target_id)
  
        if element.respond_to?('target')
          return element.target.oid
        end
  
        return element.oid
      end
  
      tailTime = getTaggedCommitTime.call(tailTag).to_i
      tailTagOid = getTagReferenceOid.call(tailTag)
      tailCommitOid = getCommitOid.call(tailTagOid)
      walker = createWalker
      
      #We need to sort topologically (since commits can be ooo)
      walker.sorting(Rugged::SORT_TOPO)
  
      if headTag.nil?
        headCommit = @repository.lookup(@repository.head.target_id)
        headTime = headCommit.time.to_i
        walker.push(headCommit.oid)
      else
        headTagOid = getTagReferenceOid.call(headTag)
        #Use same method to extract the head time.
        headTime = getTaggedCommitTime.call(headTag).to_i
        headCommitOid = getCommitOid.call(headTagOid)      
        walker.push(headCommitOid)
      end
  
      commits = Model::PACCommitCollection.new
  
      begin
        walker.each do |commit| 
          #This old method does not work as we expect. For example the ordering of the commits, which the walker
          #traverses differs from what is displayed on github. 

          if tailCommitOid != commit.oid && headTime >= commit.time.to_i && commit.time.to_i >= tailTime
            p_commit = Model::PACCommit.new(commit.oid, commit.message, commit.time)
            commits.add(p_commit)     
          end
        end
      rescue RuntimeError
        puts "Failed to load changeset"
      end
      
      return commits
    end
  end
end
