module Model

	#Model object representing a discovered task
  class PACTask
    def initialize(task_id = nil)
      #Lookup key for task management system
      #We assume that task_id is a string
      raise ArgumentError.new("PACTask init error, task id is not a string.") if not (task_id == nil or task_id.respond_to?(:to_str))
      @task_id = task_id
      #Commits tied to this task
      @commit_collection = PACCommitCollection.new
      #Data from task management systems(s)
      @attributes = { }
      #Key that determines which system(s) we need to look in for data
      @applies_to = Set.new
      #Assigned label. Used in templates so that you can group your tasks using labels.
      @label = Set.new

      # Initialized with nil - will be populated by decorators
      # FIXME is it default nil ?
      @data = nil
    end

    attr_accessor :commit_collection
    attr_accessor :task_id
    attr_accessor :data
    attr_reader   :label
    attr_reader   :applies_to
    attr_reader   :attributes

    # accepts string
    def applies_to=(val)
      raiseArgumentError.new("apllies_to method accepts only a string") if not val.is_a? String
      @applies_to << val
    end

    def add_commit(commit)
      @commit_collection.add(commit)
    end

    def commits
      @commit_collection.commits
    end

    def attributes=(value)
       @attributes = value
    end

    def to_s
      @task_id
    end

    #Apppend label to the set of existing labels. This prevents duplicate labels.
    def label=(val)
      @label << val
    end

    def clear_labels
      @label.clear
    end

    #We need to override the equals method. A task is uniquely identified by it's identifier (id). This is usually what
    #is used to fetch additonal information from the task management system
    def ==(other)
      @task_id == other.task_id
    end

    def to_liquid
      hash = {
        'task_id' => @task_id,
        'commits' => @commit_collection,
        'attributes' => attributes,
        'label' => label,
        'data' => data
      }
      hash
    end
  end

  #Model object representing a collection of tasks. Includes logic to append commits to existing tasks.
  #The task with the id 'nil' holds all the commits that do not belong to a particular task.
  class PACTaskCollection
    def initialize
      @tasks = []
    end

    #When you add a task to a collection. It will automatically add new commits to an existing task. You will
    #not get duplicate tasks in the collection.
    def add(*t)
      t.flatten.each do |task|
        unless @tasks.include? task
          @tasks << task
        else
          task.commit_collection.each do |c|
            @tasks[@tasks.index(task)].commit_collection.add(c)
          end
        end
      end
    end

    #Enumeartion method. So that PACTaskCollection.each yields a PACTask object
    def each
      @tasks.each { |task| yield task }
    end

    def length
      @tasks.length
    end

    def unreferenced_commits
      output = @tasks.select { |t| t.task_id.nil? }
      if output.first.nil?
		    return []
      else
		    return output.first.commits
      end
    end

    def referenced_tasks
      @tasks.select { |t| !t.task_id.nil? }
    end

    #Indexer method for each collected task
    def [](task_id)
      @tasks.select {|t| t.task_id == task_id}.first
    end

    #Group each task discovered by it's label
    def by_label
      labelled = Hash.new
      @tasks.each { |t|
        t.label.each { |t_label|
          unless labelled.include? t_label
            labelled[t_label] = []
          end
          unless t_label.nil?
            labelled[t_label] << t
          end
        }
      }
      labelled
    end

    def to_s
      @tasks
    end

    def to_liquid
      liquid_hash = { 'tasks' => @tasks, 'referenced' => referenced_tasks, 'unreferenced' => unreferenced_commits }
      liquid_hash.merge! by_label
      liquid_hash
    end

    attr_accessor :tasks, :attr
  end

  class PACCommitCollection
  	def initialize
  		@commits = []
  	end

  	def add(*commit)
  		@commits << commit
      @commits.flatten!
  	end

  	def each
  		@commits.each { |c| yield c }
  	end

    def count
      @commits.length
    end

    def count_with
      @commits.select{ |c| c.referenced == true }.size
    end

    def count_without
      @commits.select{ |c| c.referenced == false }.size
    end

    def health
      count_with.to_f / count
    end

  	def to_liquid
  		@commits
  	end

  	attr_accessor :commits
  end

	class PACCommit
  	def initialize(sha, message = nil, date = nil)
  		@sha = sha
  		@message = message
      @referenced = false
      @date = date
  	end

  	def to_liquid
  		{
        'sha' => @sha,
        'message' => @message,
        'header' => header,
        'referenced' => referenced,
        'shortsha' => shortsha
      }
  	end

  	def ==(other)
  		@sha == other.sha
  	end

    def header
      @message.split(/\n/).first
    end

    #Get default abbreviation from Git
    def shortsha(n = 7)
      @sha.slice(0,n)
    end

    #Match tasks agains this commit. Returns an array of matched tasks. Always contains atleast one element. Since the nil task is
    #returned for unmatch commits. That is a PACTask with an id of 'nil'.
    def matchtask(regex, split = nil)
      tasks = []
      regex.each do |r|
        @message.scan(eval(r['pattern'])).each do |arr|
          if split.nil?
            task = PACTask.new(arr[0])
            task.add_commit(self)
            self.referenced = true
            task.label = r['label']
            tasks << task
          else
            arr[0].split(split).each do |s|
              task = PACTask.new(s)
              task.add_commit(self)
              task.label = r['label']
              self.referenced = true
              tasks << task
            end
          end
        end
      end

      tasks
    end

    def to_s
      'sha:'+@sha
    end

  	attr_accessor :sha
  	attr_accessor :message
    attr_accessor :referenced
    attr_accessor :date
	end

end
