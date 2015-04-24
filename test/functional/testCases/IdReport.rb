# -*- coding: utf-8; -*-
# Test cases related to the target
#
# Copyright 2013 YXLON International A/S
#
module PAC__TestCases_IdReport
  require 'rubygems'
  gem 'test-unit'
  require 'ci/reporter/rake/test_unit_loader.rb'

  class IdReport < Test::Unit::TestCase
    require 'fileutils'
    require 'open3'
    # Order in which to run the test.
    # Order matters here as we need to run simulator first to get logfile
    # parse
    self.test_order = :alphabetic

    class << self
      # Executed before all tests
      def startup
        
      end

      # Executed after all tests
      def shutdown
      end

    end # self

    # Test setup: executed before each test
    def setup
      assert_true(system("unzip test/resources/idReportTestRepository.zip"))

    end

    # Executed after each test
    def teardown
      assert_true(system("rm -rf idReportTestRepository"))
   
    end

    # Test: short descrp
    # in depth description
    # ==== Test assumptions
    def test_runIdReport
      # Examples... should not used hardcoded commands
      assert_true(system("pwd"))
      assert_true(system("ls -al test/resources"))
      assert_true(system("./pac.rb -d 2012-01-01 --settings=test/resources/idReportTestRepository_settings.yml"))
      
    end
    

  end # class
end # module

