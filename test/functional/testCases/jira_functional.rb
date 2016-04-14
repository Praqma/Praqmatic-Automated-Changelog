# -*- coding: utf-8; -*-
module PAC__TestCases_Jira
  require 'pp'
  require 'fileutils'
  require 'ci/reporter/rake/test_unit_loader.rb'
  require_relative '../../../lib/model'
  require_relative '../../../lib/task'

  Test::Unit.at_start do
    puts "Start up jira here here"
  end

  class JiraIntegration < Test::Unit::TestCase
    require 'fileutils'
    require 'open3'

    #Example of how to construct a functional unit tests, with a known running instance of a Jira instance.
    def test_jira_case_exists
      #The only required settings to apply a task system to a named task using jira is the 'issue-link' so hotwire this, in order to get the data
      #settings = { :name => 'jira', :query_string => 'http://localhost:9090/rest/api/latest/issue/#{task_id}', :usr => 'admin', :pw => 'password' }
      #collection = Model::PACTaskCollection.new
      #task = Model::PACTask.new 'AMM-17'
      #task.applies_to = 'jira'
      #collection.add(task)
      #jira = Task::JiraTaskSystem.new(settings).apply(collection)
      #Assert that no errors returned
      #assert_true(jira)
      #Assert that the description field has a value
      #assert_not_nil(collection['AMM-17'].attributes.data.description)
    end       
  end # class
end # module
