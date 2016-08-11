require 'test/unit'
#require_relative "../../lib/model.rb"
#require_relative "../../lib/task.rb"
#require_relative "../../lib/decorators.rb"
require_relative "../../lib/core.rb"
require_relative "../../lib/logging.rb"
#require_relative "../../lib/report.rb"

class VerbosityTest < Test::Unit::TestCase

	def test_calc_verbosity_very_loud
		arguments = { '-v' => 5, '-q' => false }
		assert_equal(5, Logging.calc_verbosity(arguments))
	end

	def test_calc_verbosity_very_silent
		arguments = { '-v' => false, '-q' => 5 }
		assert_equal(-5, Logging.calc_verbosity(arguments))
	end

	def test_calc_verbosity_cancelling_values
		arguments = { '-v' => 5, '-q' => 2 }
		assert_equal(3, Logging.calc_verbosity(arguments))

		arguments = { '-v' => 1, '-q' => 1 }
		assert_equal(0, Logging.calc_verbosity(arguments))

		arguments = { '-v' => 1, '-q' => 2 }
		assert_equal(-1, Logging.calc_verbosity(arguments))
	end

	def test_verbosity_parsing_ok
		arguments = { '-v' => 5, '-q' => false }
		assert_equal(5, Logging.calc_verbosity(arguments))
		file_parsed = Core.read_settings_file(arguments)
		settings_parsed = Core.generate_settings(arguments, file_parsed)
		assert_equal(5, settings_parsed[:verbosity])		

		arguments = { '-v' => 5, '-q' => 2 }
		assert_equal(3, Logging.calc_verbosity(arguments))
		file_parsed = Core.read_settings_file(arguments)
		settings_parsed = Core.generate_settings(arguments, file_parsed)
		assert_equal(3, settings_parsed[:verbosity])		
	end

	def test_print_default
		Core.settings = { :verbosity => 0}
		assert_equal("yes", Logging.v(-1, "yes"))
		assert_equal("yes", Logging.v(0, "yes"))
		assert_nil(Logging.v(1, "no"))
	end

	def test_print_very_quiet
		Core.settings = { :verbosity => -2}
		assert_nil(Logging.v(-1, "no"))
		assert_nil(Logging.v(0, "no"))
		assert_nil(Logging.v(1, "no"))
		assert_equal("debug", Logging.v(-2, "debug"))
	end

	def test_print_very_loud
		Core.settings = { :verbosity => 2}
		assert_nil(Logging.v(3, "no"))
		assert_equal("yes", Logging.v(-3, "yes"))
		assert_equal("yes", Logging.v(0, "yes"))
		assert_equal("yes", Logging.v(1, "yes"))
		assert_equal("yes", Logging.v(2, "yes"))
	end

	def test_no_exception
		assert_nil(Logging.v(0, "no"))
		assert_nil(Logging.v(-3, "no"))
		assert_nil(Logging.v(-1, "no"))
	end

	def test_set_settings_still_logging
		assert_nil(Logging.v(0, "no"))
		Core.settings = { :verbosity => 2}
		assert_nil(Logging.v(3, "no"))
		assert_equal("yes", Logging.v(-3, "yes"))
		assert_equal("yes", Logging.v(0, "yes"))
		assert_equal("yes", Logging.v(1, "yes"))
		assert_equal("yes", Logging.v(2, "yes"))		
	end

end