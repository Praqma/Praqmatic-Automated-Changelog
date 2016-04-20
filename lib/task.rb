# encoding: utf-8
require 'pp'
require_relative 'decorators'
begin
  require 'pdfkit'
  require 'trac4r'
  require 'fogbugz'
rescue LoadError => error
  puts error
end

module Task  
  #The task system is responsible for writing the changelog. We feed it with a list of commits, and an output directory for the changelogs
  #unless otherwise specified the path will be the current directory.
  class TaskSystem
    attr_accessor :settings
    def initialize(settings)
      @settings = settings
    end

    def apply(tasks)

    end

    def html_escape_non_ascii(text)
      text.gsub(/Æ/,'&AElig;').gsub(/æ/,'&aelig;').gsub(/Ø/,'&Oslash;').gsub(/ø/,'&oslash;').gsub(/Å/,'&Aring;').gsub(/å/,'&aring;')
    end
  end

  class NoneTaskSystem < TaskSystem    
  end

  #This is the Jira task system
  class JiraTaskSystem < TaskSystem
    def initialize(settings)	
      super(settings) 
    end

    def apply(tasks)
      ok = true    
      tasks_with_no_jira_issues = []
      
      tasks.each do |t|
        begin
          if(t.applies_to.include?(@settings[:name]))  
            t.extend(JiraTaskDecorator).fetch(@settings[:query_string], @settings[:usr], @settings[:pw], debug: @settings[:debug])
          end
        rescue Exception => err   
		      tasks_with_no_jira_issues << t  
          puts "[PAC] Jira #{err.message}"
          ok = false
          t.clear_labels
          t.label = 'unknown'	
        end
      end        
      ok
    end 
  end

  class TracTaskSystem < TaskSystem
    TASK_REGEX = /Ticket\#(?<id>([0-9]+|none))+/i
    def initialize(settings)
      super(settings) 
      TracTaskDecorator.trac_instance = Trac.new settings[:trac_url], settings[:trac_usr], settings[:trac_pwd]      
    end

    def apply(tasks)
      ok = true
      tasks.each do |t|
        begin
          if(t.applies_to.include?(@settings[:name]))
            t.extend(TracTaskDecorator).fetch(debug: @settings[:debug])
          end
        rescue Exception => err
          puts "[PAC] #{err.message}"
          t.clear_labels
          t.label = 'unknown'
          ok = false
        end
      end
      ok      
    end
  end

end
