require 'test/unit'
require_relative "../../lib/model.rb"
require_relative "../../lib/task.rb"
require_relative "../../lib/decorators.rb"

class ConfigurationTest < Test::Unit::TestCase
	def test_configuration_parsing
		settings_file = File.join(File.dirname(__FILE__), '../../default_settings.yml')
		loaded = YAML::load(File.open(settings_file))
		assert_true(loaded.has_key?(:general))
		assert_equal(3, loaded[:task_systems].length)
		assert_equal(3, loaded[:templates].length)		
	end
end