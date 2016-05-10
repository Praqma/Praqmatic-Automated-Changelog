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
		settings_parsed = Core.generate_settings(arguments)
		assert_equal('PAC Changelog Name Override', settings_parsed[:properties]['title'])		
		defined = Report::Generator.new.define_properties(settings_parsed)
		assert_equal('PAC Changelog Name Override', defined['properties']['title'] )
		assert_equal('PAC Changelog Name Override', settings_parsed[:properties]['title'])
	end

	#The everything ok scenario where properties is not specified. Assert that the default 'title' is present 
	def test_properties_parse_omitted
		settings_file = File.join(File.dirname(__FILE__), '../../default_settings.yml')
		arguments = { '--settings' => "#{settings_file}" }	
		settings_parsed = Core.generate_settings(arguments)
		defined = Report::Generator.new.define_properties(settings_parsed)
		assert_equal('PAC Changelog', defined['properties']['title'] )
		assert_nil(settings_parsed[:properties])
	end

	#The sceenario where the JSON is invalid and cannot be parsed. Assert that an exception is thrown
	def test_properties_incorrect_json
		settings_file = File.join(File.dirname(__FILE__), '../../default_settings.yml')
		arguments = { '--settings' => "#{settings_file}", '--properties' => "{ title'PAC Chang} " }
		assert_raise do |err|
			settings_parsed = Core.generate_settings(arguments)
		end
	end
end