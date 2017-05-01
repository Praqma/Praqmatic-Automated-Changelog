require 'test/unit'
require_relative "../../lib/model.rb"
require_relative "../../lib/task.rb"
require_relative "../../lib/decorators.rb"

class ModelTest < Test::Unit::TestCase


	# Testing new PAC task creation
	# but not the trivial get/setters
	def test_PACTask_initialize

		# A PAC task initializes with a task_id, either nil or something that makese
		# sense later in the system. This means a string.
		tc = Model::PACTask.new()
		assert_nil(tc.task_id,"Task initialized without task id should task_id nil")

		# Shouldn't reaise argument error as we expect nil or string
		assert_nothing_raised(ArgumentError) {
			tc = Model::PACTask.new('issuekey')
			tc = Model::PACTask.new()
		}
	end


	# Applies to method is setting the tasks systems the PAC tasks have data in
	def test_PACtask_applies_to_method
		ts = Model::PACTask.new()


		assert_equal(Set.new,ts.applies_to(),"Applies to is default empty set.")
		# Applies to take a string matching the entry for task system in the config file
		ts.applies_to='mytasksystem'
		tasksystem = Set.new ['mytasksystem']
		assert_equal(tasksystem,ts.applies_to(),"Applies to tasks system not 'mytasksystem' as expected.")
	end

	def test_model_task_collection
		tc = Model::PACTaskCollection.new
		arr_of_tasks = [ Model::PACTask.new("1"), Model::PACTask.new("2")]
		tc.add(arr_of_tasks)

		ta = Model::PACTask.new("3")
		tc.add(ta)
		tc.add(ta)
		assert_true(tc.length == 3)
	end

	def test_commit_collection_add
		cc = Model::PACCommitCollection.new
		commit1 = Model::PACCommit.new('abcd1234')
		commit1.referenced = true
		commit2 = Model::PACCommit.new('abcd1234abcd')
		commit3 = Model::PACCommit.new('abcd12341235')
		cc.add(commit1)
		cc.add(commit2,commit3)

		assert_true(cc.count_with  == 1)
		assert_true(cc.count_without == 2, "Commit list must contain 2 unreferenced commits, it contained #{cc.count_without}")
		assert_equal(33, (cc.health*100).to_i, "We expect the health to be 33, it was #{(cc.health*100).to_i}")
	end

	def test_commit_associated_with_more_than_one_case
		tc = Model::PACTaskCollection.new
		commit1 = Model::PACCommit.new('abcd1234')
		commit2 = Model::PACCommit.new('abcd1234abcd')
		commit3 = Model::PACCommit.new('ffff88888')

		task1 = Model::PACTask.new('task1')
		task1.add_commit(commit1)

		task11 = Model::PACTask.new('task1')
		task11.add_commit(commit2)

		task2 = Model::PACTask.new('task2')
		task2.add_commit(commit1)

		task3 = Model::PACTask.new('task3')
		task3.add_commit(commit3)

		arr_task = [task1, task2, task3]
		tc.add(arr_task)

		assert_true(tc.length == 3, "Task list must contain 3 items, it contained #{tc.length}")

		assert_not_nil(tc['task1'])
		assert_not_nil(tc['task2'])
		assert_not_nil(tc['task3'])
		assert_nil(tc['task4'])
	end

	def test_decorator_module
		task11 = Model::PACTask.new('STACI-2')
		#The collection of tasks is totally decoupled from the task system itself. Only requirement is that
		#you should be able to 'fetch' the required data from a task, given the id of the object. How that is done
		#is an implementaion detail.
		task_c = Model::PACTaskCollection.new
		task_c.add(task11)
		assert_false(task11.attributes.has_key?('data'))
		task11.extend(JiraTaskDecorator)
		#Notice how that 'task11' now responds to the attributes call with the data field (empty at the moment)
		puts task11.attributes
		assert_true(task11.attributes.has_key?('data'))
	end


	def test_grouping
		tc_grouped = Model::PACTaskCollection.new
		commit1 = Model::PACCommit.new('abcd1234')
		group1 = [Model::PACTask.new('STACI-3A'), Model::PACTask.new('STACI-3B'), Model::PACTask.new('STACI-3C')]

		group1.each do |t|
			t.label = 'STACI'
			t.add_commit(commit1)
			t.applies_to = 'jira1'
			t.applies_to = 'jira2'
		end

		commit2 = Model::PACCommit.new('abcd1234abcd')
		group2 = Model::PACTask.new('STACI-4')
		group2.label = 'LUCI'


		group2.add_commit(commit2)

		tc_grouped.add(group1)
		tc_grouped.add(group2)

		assert_equal(3, tc_grouped.by_label["STACI"].length)
		assert_equal(1, tc_grouped.by_label["LUCI"].length)

		group1.each do |g1|
			assert_equal(2, g1.applies_to.length)
		end
	end

end
