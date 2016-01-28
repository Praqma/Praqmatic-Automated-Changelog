# encoding: utf-8
require 'yaml'
require_relative "./task.rb"
require_relative "./gitvcs"
require_relative "./mercurialvcs"

module Core extend self

  def settings
    @@settings
  end

  def settings=(val)
    @@settings = val
  end
  
  #Requires a configuration section for the task system to be applied
  def apply_task_system(task_system, tasks)
    if task_system[:name] == 'trac'
      Task::TracTaskSystem.new(@@settings).apply(tasks)
    end
    if task_system[:name] == 'jira'
      Task::JiraTaskSystem.new(@@settings).apply(tasks)
    end   
  end
  
  def vcs
    if @@settings[:vcs][:type] == 'git'
      Vcs::GitVcs.new(settings[:vcs])
    elsif @@settings[:vcs][:type] == 'hg'
      Vcs::MercurialVcs.new(@@settings[:vcs])
    else
      raise ArgumentError, 'The configuration settings does not include any supported (d)vcs'
    end
  end
  
  def to_time(datestring)
    DateTime.strptime(datestring, @settings[:general]['date_template']).to_time    
  end

  #This is now core functionality. The task of generating a collection of tasks based on the commits found
  #This takes in a PACCommitCollection and returns a PACTaskCollection 
  def task_id_list(commits)
    regex_arr = []

    tasks = Model::PACTaskCollection.new

    commits.each do |c_pac|

      referenced = false
      #Regex ~ Eacb regex in the task system
      settings[:task_systems].each do |ts|
        #Loop over each task system. Parse commits for matches
        if ts.has_key? :delimiter 
          split_pattern = eval(ts[:delimiter]) 
        end

        if ts.has_key? :regex 
          tasks_for_commit = c_pac.matchtask(ts[:regex], split_pattern)
          tasks_for_commit.each do |t|
            t.applies_to = ts[:name]
          end            
          #If a task was found
          unless tasks_for_commit.empty?
            referenced = true          
            tasks.add(tasks_for_commit)                
          end          
        end
      end

      if !referenced
        task = Model::PACTask.new
        task.add_commit(c_pac)
        tasks.add(task)
      end      

    end

    tasks      
  end

end
