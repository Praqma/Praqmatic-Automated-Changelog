# encoding: utf-8
require 'pp'
require_relative 'decorators'

module Task
  #The task system is responsible for writing the changelog. We feed it with a list of commits, and an output directory for the changelogs
  #unless otherwise specified the path will be the current directory.
  class TaskSystem
    attr_accessor :settings
    def initialize(settings)
      @settings = settings
    end

    def apply(tasks)
      true
    end
  end

  class NoneTaskSystem < TaskSystem
  end

  #This is the Jira task system
  class JsonTaskSystem < TaskSystem
    attr_accessor :general_settings

    def initialize(settings)
      super(settings)
    end

    def apply(tasks)
      ok = true
      tasks_without_issues = []

      tasks.each do |t|
        begin
          if(t.applies_to.include?(@settings[:name]))
            t.extend(JsonTaskDecorator).fetch(@settings[:query_string], @settings[:auth], Core.settings[:general][:ssl_verify])
            Logging.verboseprint(1, "[PAC] Applied task system Json to #{t.task_id}")
          end
        #This handles the case where we matched the regex. But the user might have a typo in the issue id.
        #This means the issue cannot be looked up.
        rescue Exception => err
		      tasks_without_issues << t
          Logging.verboseprint(0, "[PAC] Json #{err.message}")
          Logging.verboseprint(1, err.backtrace)
          ok = false
          t.clear_labels
          t.label = 'unknown'
        end
      end
      ok
    end
  end
end
