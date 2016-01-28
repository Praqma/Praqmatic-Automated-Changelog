require 'test/unit'
require_relative "../../lib/model.rb"
require_relative "../../lib/task.rb"
require_relative "../../lib/decorators.rb"

class JiraTaskSystemTest < Test::Unit::TestCase
	#This tests that when you query a an adress that does not exist, we exhibit appropriate behaviour and handle it gracefully. This is a test for a domain that 
	#actually exits, the server will then throw a 404 not found.  
	def test_jira_rainy_day_scenario_host_exits_page_not_found
		settings = { :name => 'jira', 'query_string' => 'http://www.praqma.com/#{task_id}', 'usr' => 'usr', 'pw' => 'pw' }
		collection = Model::PACTaskCollection.new
		task = Model::PACTask.new 'Jira1'
		task.applies_to = 'jira'
		collection.add(task)
		jira = Task::JiraTaskSystem.new(settings).apply(collection)
		assert_false(jira)
	end

	#This tests that when you query a an adress that does not exist, we exhibit appropriate behaviour and handle it gracefully. This is a test a site that
	#does not exists.
	def test_jira_rainy_day_scenario_host_not_exists
		settings = { :name => 'jira', 'query_string' => 'http://we.are.not.doughtnuts.fix/#{task_id}', 'usr' => 'usr', 'pw' => 'pw' }
		collection = Model::PACTaskCollection.new
		task = Model::PACTask.new 'Jira1'
		task.applies_to = 'jira'
		collection.add(task)
		jira = Task::JiraTaskSystem.new(settings).apply(collection)
		assert_false(jira)
	end		

	#Basic test to verify the correct response  
	def test_returned_incomplete_json
		settings = { :name => 'jira', 'query_string' => 'http://we.are.not.doughtnuts.fix/#{task_id}', 'usr' => 'usr', 'pw' => 'pw' }
		collection = Model::PACTaskCollection.new
		task = Model::PACTask.new 'Jira1'
		task.applies_to = 'jira'
		collection.add(task)
		jira = Task::JiraTaskSystem.new(settings).apply(collection)

		assert_raise JSON::ParserError do |err|
			task.parse('incomplete_json {')
		end
		
		assert_false(jira)		
	end
end