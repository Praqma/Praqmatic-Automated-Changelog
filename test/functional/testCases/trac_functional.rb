# -*- coding: utf-8; -*-
module PAC__TestCases_Trac
  require 'pp'
  require 'fileutils'
  require 'zip/zip'
  require 'ci/reporter/rake/test_unit_loader.rb'
  require_relative '../../../lib/model'
  require_relative '../../../lib/decorators'
  require_relative '../../../lib/task'

  class TracIntegration < Test::Unit::TestCase
    require 'fileutils'
    require 'open3'

    Test::Unit.at_start do
      #Execute the startup script for Trac
      trac_host_port = ENV['HOST_PORT'] || '28080'

      puts %x( ./test/resources/start_task_system.sh "trac" #{trac_host_port} )
      
      settings = { :trac_url => "http://localhost:#{trac_host_port}/trac", :trac_usr => 'admin', :trac_pwd => 'UqmvS76r7D', :name => 'trac' }
      system = Task::TracTaskSystem.new(settings)
      TracTaskDecorator.trac_instance.tickets.create "First ticket", "First ticket description", :type => 'defect', :version => '1.0', :milestone => 'poc'
      TracTaskDecorator.trac_instance.tickets.create "Second ticket", "Second ticket description", :type => 'defect', :version => '1.0', :milestone => 'poc'
      TracTaskDecorator.trac_instance.tickets.create "Third ticket", "Third ticket description", :type => 'defect', :version => '1.0', :milestone => 'poc'
    end

    Test::Unit.at_exit do
      bn = ENV['BUILD_NUMBER'] || '0000'
      puts %x( ./test/resources/stop_task_system-trac-#{bn}.sh )
    end

    def test_http_not_ok
      trac_host_port = ENV['HOST_PORT'] || '28080'
      #The only required settings to apply a task system to a named task using jira is the 'issue-link' so hotwire this, in order to get the data
      #settings[:trac_url], settings[:trac_usr], settings[:trac_pwd]
      settings = { :trac_url => "http://localhost:#{trac_host_port}/trac", :trac_usr => 'admin', :trac_pwd => 'UqmvS76r7D', :name => 'trac' }
      system = Task::TracTaskSystem.new(settings)    
      collection = Model::PACTaskCollection.new      
      task = Model::PACTask.new 666.to_s
      task.applies_to = 'trac'
      task.label = 'found'
      collection.add(task)

      trac = system.apply(collection)
      #Assert that errors errors returned
      assert_false(trac)
      #We assign the label 'unknown' to tasks that failed to fetch metadata'
      assert_true(task.label.include?('unknown'))
      #Also assert that the label 'found' has been cleared
      assert_false(task.label.include?('found'))
    end

    #First functional 
    def test_trac_case_exists
      trac_host_port = ENV['HOST_PORT'] || '28080'
      #The only required settings to apply a task system to a named task using jira is the 'issue-link' so hotwire this, in order to get the data
      #settings[:trac_url], settings[:trac_usr], settings[:trac_pwd]
      settings = { :trac_url => "http://localhost:#{trac_host_port}/trac", :trac_usr => 'admin', :trac_pwd => 'UqmvS76r7D', :name => 'trac' }
      system = Task::TracTaskSystem.new(settings)    

      collection = Model::PACTaskCollection.new

      [1,2,3].each do |ticket| 
        task = Model::PACTask.new ticket.to_s
        task.applies_to = 'trac'
        collection.add(task)
      end 
      
      trac = system.apply(collection)
      #Assert that no errors returned
      assert_true(trac)

      assert_equal("First ticket description", collection['1'].attributes['data'][:description])
      assert_equal("Second ticket description", collection['2'].attributes['data'][:description])
      assert_equal("Third ticket description", collection['3'].attributes['data'][:description])

      assert_equal("First ticket", collection['1'].attributes['data'][:summary])
      assert_equal("Second ticket", collection['2'].attributes['data'][:summary])
      assert_equal("Third ticket", collection['3'].attributes['data'][:summary])
    end       
  end # class
end # module