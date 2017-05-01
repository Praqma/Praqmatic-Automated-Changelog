require_relative '../lib/model.rb'

RSpec.describe Model do
	describe "class PACTask" do
		describe "method: new" do
			context "When a PACTask is created" do
				it "the data object should not be there" do
					ta = Model::PACTask.new("1")
					expect(ta.data).to be_nil
				end
			end
		end
		describe "applies_to" do
			let(:ta) { Model::PACTask.new("2") }
			let(:ta_full) {
				m = Model::PACTask.new("2")
				m.applies_to = 'jira'
				m
			}

			context "Given a new task" do
				it "applies_to should be empty" do
					expect(ta.applies_to.length).to be 0
				end

				it "should append the applied string when called" do
					ta.applies_to = 'jira'
					expect(ta.applies_to.length).to be 1
				end

				it "should append the new applied string to the end of the list when re-applied" do
					ta_full.applies_to = 'trac'
					expect(ta_full.applies_to.length).to be 2
				end

			end
		end

		describe "add_label" do
			let(:std_task) { Model::PACTask.new("1") }
			context "Given a new task" do
				it "label should be empty" do
					expect(std_task.label.length).to be 0
				end
			end

			context "When the user adds a label to the task" do
				it "the label should be added to the list" do
				end
			end
		end
	end
	describe "class: PACTaskCollection" do
		context "Given a new collection" do
			it "HAHA" do
				expect(1).to be 1
			end
		end
	end
end
