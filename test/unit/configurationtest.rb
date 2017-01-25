require 'test/unit'
require_relative "../../lib/model.rb"
require_relative "../../lib/task.rb"
require_relative "../../lib/decorators.rb"
require_relative "../../lib/core.rb"
require_relative "../../lib/report.rb"

class ConfigurationTest < Test::Unit::TestCase

	def test_configuration_parsing
		settings_file = File.join(File.dirname(__FILE__), '../../default_settings.yml')
		loaded = YAML::load(File.open(settings_file))
		assert_true(loaded.has_key?(:general))
		assert_equal(3, loaded[:task_systems].length)
		assert_equal(3, loaded[:templates].length)		
	end

	#The everything ok scenario
	def test_properties_parsing_ok
		settings_file = File.join(File.dirname(__FILE__), '../../default_settings.yml')
		arguments = { '--settings' => "#{settings_file}", '--properties' => '{"title" : "PAC Changelog Name Override" }' }
		file_parsed = Core.read_settings_file(arguments)
		settings_parsed = Core.generate_settings(arguments, file_parsed)
		assert_equal('PAC Changelog Name Override', settings_parsed[:properties]['title'])		
		defined = Report::Generator.new.to_liquid_properties(settings_parsed)
		assert_equal('PAC Changelog Name Override', defined['properties']['title'] )
		assert_equal('PAC Changelog Name Override', settings_parsed[:properties]['title'])
	end

	#The sceenario where the JSON is invalid and cannot be parsed. Assert that an exception is thrown
	def test_properties_incorrect_json
		settings_file = File.join(File.dirname(__FILE__), '../../default_settings.yml')
		arguments = { '--settings' => "#{settings_file}", '--properties' => "{ title'PAC Chang} " }
		file_parsed = Core.read_settings_file(arguments)		
		assert_raise do |err|
			settings_parsed = Core.generate_settings(arguments, file_parsed)
		end
	end
	
	#Credentials test (test the -c option. for username and password overrides)
	def test_configure_credentials
		settings_file = File.join(File.dirname(__FILE__), '../../default_settings.yml')
		#Notice the wierd way docopt handles it. The -c flag is a repeat flag, each option is then grouped positionally. So for each 'c' specified 
		#c is incremented, and the index of the then the value specified.
		arguments = { '--settings' => "#{settings_file}", '--properties' => '{"title" : "PAC Changelog Name Override" }', '-c' => 2, 
			'<user>' => ["newuser", "tracuser"], 
			'<password>' => ["newpassword", "tracpassword"], 
			'<target>' => ["jira", "trac"] } 

		file_parsed = Core.read_settings_file(arguments)
		settings_parsed = Core.generate_settings(arguments, file_parsed)		
		assert_equal('newuser', settings_parsed[:task_systems][1][:usr])
		assert_equal('newpassword', settings_parsed[:task_systems][1][:pw])

		assert_equal('tracuser', settings_parsed[:task_systems][2][:usr])
		assert_equal('tracpassword', settings_parsed[:task_systems][2][:pw])
	end

	def test_raise_exception_on_missing_settings_file 
		config = { '--settings' => 'not-there.yml' }
		assert_raise RuntimeError do 
			Core.read_settings_file(config)
		end
		#Assert the inverse as well
		config['--settings'] = 'pac_settings.yml'
		Core.read_settings_file(config)
	end	

end