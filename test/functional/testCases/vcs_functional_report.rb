# -*- coding: utf-8; -*-
module PAC__TestCases_IdReport
  require 'pp'
  require 'fileutils'
  require 'zip/zip'
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

      Core.settings = settings
      #PACCommitCollection
      commit_map = Core.vcs.get_delta("f9a66ca6d2e616b1012a1bdeb13f924c1bc9b4b6", "fb493078d9f42d79ea0e3a56abca7956a0d47123")
      pp "#################################"
      pp "Commit map for the repository is:"
      pp commit_map

      pp "##############################################"
      pp "grouped by tasks ids the commits are"
      #PACTaskCollection
      task_list = Core.task_id_list(commit_map)
      pp task_list

      pp "Checking with test asserts if the id expected are there:"
      # based on our created test repository
      assert_false(task_list["1"].nil?, "Didn't find task reference for '1' in the repository as expected")
      assert_false(task_list["2"].nil?, "Didn't find task reference for '2' in the repository as expected #{task_list.tasks}")
      assert_false(task_list["3"].nil?, "Didn't find task reference for '3' in the repository as expected")
      pp "DONE - asserts okay for expected ids"

      pp "##############################################"
      pp "Un-referenced commits are:"
      pp task_list.unreferenced_commits

      pp "Checking with test asserts for un-referenced commits :"
      # Checking on unreferenced shas:
      pp task_list.unreferenced_commits
      #assert_equal(task_list.unreferenced.first.commits.class, Array)
      #assert_equal(2, task_list.unreferenced.first.commit_collection.commits.length)
      assert_true(task_list["none"].commits.include?(Model::PACCommit.new("a789b472150f462a8ae291577dcf7557b2b4ca55")))
      assert_true(task_list.unreferenced_commits.include?(Model::PACCommit.new("55857d4e9838d1855b10e4c30b43a433e2db47cd")))
      pp "DONE - assert okay for un-referenced commits"

    end

  end # class
end # module

