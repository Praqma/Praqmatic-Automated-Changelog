# -*- coding: utf-8; -*-
# The PAC test cases for jira uses the setup described and documented under the 
# test/resources/jira-env folder of this project. It uses a preconfigured Jira container with pre-made
# issues. The pre-filled data (json) is included in the folder along with instructions on how to add this data.
# The functional tests in this module make uses of that.  
module PAC__TestCases_Jira
  require 'pp'
  require 'fileutils'
  require 'ci/reporter/rake/test_unit_loader.rb'
  require_relative '../../../lib/model'
  require_relative '../../../lib/task'

  Test::Unit.at_start do
    #Execute the startup script for Jira.
    jira_host_port = ENV['HOST_PORT'] || '28080'
    puts %x( ./test/resources/start_task_system.sh "jira" #{jira_host_port} )
  end

  Test::Unit.at_exit do
    bn = ENV['BUILD_NUMBER'] || '0000'
    puts %x( ./test/resources/stop_task_system-jira-#{bn}.sh )
  end  

  class JiraIntegration < Test::Unit::TestCase
    require 'fileutils'
    require 'open3'

    #Inspect one case throughly. In this case it's the FAS-1 case. See the test/resources/jira-env/FAS-1.json file on how Jira json output looks like.
    def test_jira_case_exists
      jira_host_port = ENV['HOST_PORT'] || '28080'
      #The only required settings to apply a task system to a named task using jira is the 'issue-link' so hotwire this, in order to get the data

      #We split the string into the static part here
      static = "http://localhost:#{jira_host_port}"
      #The query string is construted this way so that the 'task_id' is only evaluated when jira is 'applied' to the task. Therefore we have
      #to wrap that part in single plings.
      settings = { :name => 'jira', :query_string => static+'/rest/api/latest/issue/#{task_id}', :usr => 'admin', :pw => 'admin' }
      collection = Model::PACTaskCollection.new
      task = Model::PACTask.new 'FAS-1'
      task.applies_to = 'jira'
      collection.add(task)

      #First assert that the inverse is true before starting (the attributes field SHOULD be empty here)
      assert_true(task.attributes.empty?)

      jira = Task::JiraTaskSystem.new(settings).apply(collection)

      #Assert data was there
      assert_false(task.attributes.empty?)

      #Verify that the data is correct.
      description = task.attributes['data']['fields']['description']
      summary = task.attributes['data']['fields']['summary']
      assert_equal("Thatcher's government since 1952). Scotland The Scots pine marten.", summary)
      assert_true(description.include?("After 1860 at Holyrood in the Great Highland English"))
    end       

    #This test just processes a larger collection of issues, and makes sure that data is there. Data is not verified, we're just testing that the
    #'attributes' of each case has been set.
    def test_bulk_procesing_case
      jira_host_port = ENV['HOST_PORT'] || '28080'
      
      #We split the string into the static part here
      static = "http://localhost:#{jira_host_port}"

      #The query string is construted this way so that the 'task_id' is only evaluated when jira is 'applied' to the task. Therefore we have
      #to wrap that part in single plings.
      settings = { :name => 'jira', :query_string => static+'/rest/api/latest/issue/#{task_id}', :usr => 'admin', :pw => 'admin' }
      collection = Model::PACTaskCollection.new
      (2..30).each do |tid|
        task = Model::PACTask.new "FAS-#{tid}"
        task.applies_to = 'jira'
        collection.add(task)        
      end

      #Assert that each task has some assigned attributes
      jira = Task::JiraTaskSystem.new(settings).apply(collection)
      collection.each do |task|
        assert_false(task.attributes.empty?)
      end       
    end

    #This tests the case where the developer has put in the wrong task number (one that doesn't exists)
    #The same response is recorded for ALL return codes that are not 200 OK (unless res.is_a? Net::HTTPOK -This includes 400 Bad Request, 403 Permission denied etc...)
    #TODO: How do we handle and test http-redirects?
    def test_http_not_ok
      jira_host_port = ENV['HOST_PORT'] || '28080'
      
      #We split the string into the static part here
      static = "http://localhost:#{jira_host_port}"

      #The query string is construted this way so that the 'task_id' is only evaluated when jira is 'applied' to the task. Therefore we have
      #to wrap that part in single plings.
      settings = { :name => 'jira', :query_string => static+'/rest/api/latest/issue/#{task_id}', :usr => 'admin', :pw => 'admin', :debug => true}

      #Create a non-exiting issue that was referenced
      collection = Model::PACTaskCollection.new
      task = Model::PACTask.new 'FAS-90029'
      task.applies_to = 'jira'
      task.label = 'found'
      collection.add(task)

      #The correct behaviour in this case is that 'jira' becomes false, an error messsage is written to std.out 
      jira = Task::JiraTaskSystem.new(settings).apply(collection)
      assert_false(jira)

      #We assign the label 'unknown' to tasks that failed to fetch metadata. (The commit is STILL referenced because it matched our regex)
      assert_true(task.label.include?('unknown'))
      #Also assert that the label 'found' has been cleared
      assert_false(task.label.include?('found'))
    end

  end # class
end # module
