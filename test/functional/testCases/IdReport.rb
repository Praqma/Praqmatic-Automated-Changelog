# -*- coding: utf-8; -*-
# Test cases related to the target
#
# Copyright 2013 YXLON International A/S
#
module PAC__TestCases_IdReport
  require 'pp'
  require 'fileutils'
  gem 'test-unit'
  require 'zip/zip'
  require 'test/unit'
  require 'ci/reporter/rake/test_unit_loader.rb'

  # The idReport test cases tests our id report functionality
  # by using a git repository created for the tests with known
  # content like SHAs, commits etc.
  # As the not all source code is written easy to test, we call into
  # a few functions and test on their return values.
  class IdReport < Test::Unit::TestCase
    require 'fileutils'
    require 'open3'
    # Order in which to run the test.
    # Order matters here as we need to run simulator first to get logfile
    # parse
    self.test_order = :alphabetic

    # Helper function to unzip our repo
    def unzip_file (file, destination)
      Zip::ZipFile.open(file) { |zip_file|
        zip_file.each { |f|
          f_path=File.join(destination, f.name)
          FileUtils.mkdir_p(File.dirname(f_path))
          zip_file.extract(f, f_path) unless File.exist?(f_path)
        }
        puts "Successfully extracted"
      }
    end

    # Test setup: executed before each test
    def setup
      unzip_file("test/resources/idReportTestRepository.zip", "test/resources")
    end

    # Executed after each test
    def teardown
      FileUtils.rm_rf("test/resources/idReportTestRepository")
    end

    def test_mytest
      require_relative '../../../lib/core'
      settings_file = File.join(File.dirname(__FILE__), '../../resources/idReportTestRepository_settings.yml')

      settings = YAML::load(File.open(settings_file))

      Core.load(settings)

      commit_map = Core.vcs.get_commit_messages_by_commit_sha("f9a66ca6d2e616b1012a1bdeb13f924c1bc9b4b6", "fb493078d9f42d79ea0e3a56abca7956a0d47123")
      pp "#################################"
      pp "Commit map for the repository is:"
      pp commit_map

      pp "##############################################"
      pp "grouped by tasks ids the commits are"
      grouped_by_task_id = Core.task_system.task_id_list(commit_map)
      pp grouped_by_task_id

      pp "Checking with test asserts if the id expected are there:"
      # based on our created test repository
      assert_true(grouped_by_task_id.has_key?("1"), "Didn't find task reference for '1' in the repository as expected")
      assert_true(grouped_by_task_id.has_key?("2"), "Didn't find task reference for '2' in the repository as expected")
      assert_true(grouped_by_task_id.has_key?("3"), "Didn't find task reference for '3' in the repository as expected")
      pp "DONE - asserts okay for expected ids"

      pp "##############################################"
      pp "Un-referenced commits are:"
      unreferenced = Core.task_system.get_shas_without_reference(commit_map, grouped_by_task_id)
      pp unreferenced

      pp "Checking with test asserts for un-referenced commits :"
      # Checking on unreferenced shas:
      assert_true(unreferenced.include?("f9a66ca6d2e616b1012a1bdeb13f924c1bc9b4b6"))
      assert_true(unreferenced.include?("a789b472150f462a8ae291577dcf7557b2b4ca55"))
      assert_true(unreferenced.include?("55857d4e9838d1855b10e4c30b43a433e2db47cd"))
      pp "DONE - assert okay for un-referenced commits"

    end

  end # class
end # module

