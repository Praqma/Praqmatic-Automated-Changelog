require 'test/unit'
require_relative "../../lib/model.rb"
require_relative "../../lib/task.rb"
require_relative "../../lib/decorators.rb"

class FogBugzTaskSystemTest < Test::Unit::TestCase

  Test::Unit.at_start do  
  	Core.settings = { :verbosity => 1}
  end

  #Just testing that we fail to apply the task system when there is no task system 
	def test_fogbugz_configuration_should_fail
		settings = { :name => 'fogbugz', :query_string => 'http://www.praqma.com/#{task_id}', :usr => 'usr_fb', :pw => 'pw_fb' }
		collection = Model::PACTaskCollection.new
		task = Model::PACTask.new 'Fogbugz1'
		task.applies_to = 'fogbugz'
		collection.add(task)
		fogbugz = Task::FogBugzTaskSystem.new(settings).apply(collection)
		assert_false(fogbugz)
	end

end

