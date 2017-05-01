require_relative '../lib/core'
require 'docopt'

RSpec.describe "Command line interactions" do
	describe "Commands that must be checked" do

		context "When the user wants to get help with the --help option" do
			it "prints out the help message on the screen" do
				ARGV = "--help" 
				expect { Docopt::docopt(Core.cli_text(__FILE__)) }.to raise_error(Docopt::Exit)
 			end
 		end

 		context "When an illegalt parameter is passed" do
 			it "prints the help message to the screen and informs the user" do
				ARGV = "--mads" 
				#TODO: Check error message
				expect { Docopt::docopt(Core.cli_text(__FILE__)) }.to raise_error(Docopt::Exit)
 			end
 		end

 		context "When from-latest-tag is used" do
 			it "does as we expect" do
				ARGV = "from-latest-tag '*' --to='somesha' --settings='pac_settings.yml'" 
				#TODO: We should add some kind of validation test here
				input = Docopt::docopt(Core.cli_text(__FILE__)) 				
 			end
		end

 		context "When the command from is used" do
 			it "does as we expect" do
				ARGV = "from 'sha' --to='somesha' --settings='pac_settings.yml' -c 'user' 'pass' 'tasksystem'" 
				input = Docopt::docopt(Core.cli_text(__FILE__))
				puts input 				
 			end
		end	
	end
end