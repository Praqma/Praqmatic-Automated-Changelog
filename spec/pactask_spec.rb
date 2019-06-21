require_relative '../lib/model.rb'

RSpec.describe Model do
	describe "class PACTask" do

		describe ".new" do
			context "Task initialized" do
			  let(:pt) { Model::PACTask.new() }
			  it "without task id should have task_id nil" do
				  expect(pt.task_id).to eq(nil)
			  end

				let(:pt_id) { Model::PACTask.new(3) }
				it "with something else than a string should raiseArgument error as we expect it to be a string later" do
					expect{Model::PACTask.new(3)}.to raise_error(ArgumentError)
				end
			end
		end

		describe "applies_to" do
			let(:pt) { Model::PACTask.new() }

			context "Given a new task" do
				it "applies_to should be empty" do
					expect(pt.applies_to.length).to be 0
				end

				it "applies_to should return the set of applied tasks systems" do
					pt.applies_to = 'mytaskssystem'
					ts = Set.new ['mytaskssystem']
					expect(pt.applies_to).to eq(ts) # not object identity, just content
				end
			end
		end
	end
end
