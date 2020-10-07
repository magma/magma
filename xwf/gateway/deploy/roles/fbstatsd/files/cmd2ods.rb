#!/usr/bin/env ruby
# A script that pasrses yml files containing comand_line and cron spec.
# Runs the commands and sends the results to ODS via fbstatsd.
# This assumes fbstatsd daemon is running on the box
# Example yml config file:
## policies:
##   'command_line''' "ip x p c | awk '{print $3}'"
##   'cron':
##     'minute': 1

require 'yaml'
require 'time'
require 'socket'
require 'logger'
require 'timeout'

# rubocop:disable all

CONFIG_DIR = File.join('/', 'opt', 'xwf', 'cmd2ods', '*.yml')
DEFAULT_TIMEOUT = 60

$logger = Logger.new(STDOUT)
$logger.level = Logger::INFO

$config = { :server => 'localhost', :port => 1055 }

class Command
  include Socket::Constants

  attr_accessor :execution_time_in_sec
  attr_accessor :command_line
  attr_accessor :ods_key
  attr_accessor :timeout

  def initialize(ods_key, command_line, cron, timeout = DEFAULT_TIMEOUT)
    self.command_line = command_line
    self.ods_key = ods_key
    self.timeout = timeout
    second = cron['second'].nil? ? 0 : cron['second']
    minute = cron['minute'].nil? ? 0 : cron['minute']
    hour = cron['hour'].nil? ? 0 : cron['hour']
    day = cron['day'].nil? ? 0 : cron['day']
    self.execution_time_in_sec = calculate_execution_time(
      second, minute, hour, day
    )
    $logger.info "Initialized Command #{self.ods_key} every \
                        #{execution_time_in_sec} seconds"
  end

  def calculate_execution_time(second = 0, minute = 0, hour = 0, day = 0)
    second + (minute * 60) + (hour * 60 * 60) + (day * 60 * 60 * 24)
  end

  # Run the command, returns output
  def execute
    begin
      rout, wout = IO.pipe
      rerr, werr = IO.pipe
      stdout, stderr = nil

      pid = Process.spawn(command_line, :pgroup => true, :out => wout, :err => werr)
      Timeout.timeout(timeout) do
        Process.wait(pid)

        wout.close
        werr.close

        stdout = rout.readlines.join
        stderr = rerr.readlines.join
      end
    rescue Timeout::Error
      Process.kill(-9, pid)
      Process.detach(pid)
    ensure
      wout.close unless wout.closed?
      werr.close unless werr.closed?

      rout.close
      rerr.close
    end
    stdout
  end

  # This is the format that fbstatsd is expecting
  def format_ods_string(key, value)
    "#{key}:#{value}|g"
  end

  # Sends a string to udp socket
  def send_to_ods(ods_string)
    sock = UDPSocket.new
    sock.send(ods_string, 0, $config[:host], $config[:port])
    $logger.info "Wrote to ODS: #{ods_string}"
  end

  # Run the command, format the output and send it to ODS
  def run
    result = execute
    result.split("\n").each do |r|
      parsed = r.split
      ods_value = parsed[-1]
      line_key = parsed[0..-2].join('_')
      if ods_value !~ /^[-+]?[0-9]*\.?[0-9]+$/
        $logger.warn "Command #{ods_key} returned non numeric value - \
                  \"#{ods_value}\", when running \"#{command_line}\""
        return
      end
      ods_key = line_key.empty? ? self.ods_key : [self.ods_key, line_key].join('.')
      ods_string = format_ods_string(ods_key, ods_value)
      send_to_ods(ods_string)
    end
  rescue Exception => ex
    $logger.error "Could not run #{self.ods_key}"
    $logger.error ex
    ex.backtrace.each do |b|
      $logger.error "\t#{b}"
    end
  end
end

# Read and join all yml files in to a hash
commands_configs = {}
Dir[CONFIG_DIR].each do |config_file|
  # rubocop:disable Security/YAMLLoad
  commands_configs.merge! YAML.load(File.read(config_file))
end

# Array that holds all commands as Command class objects
commands = []

commands_configs.keys.each do |command_config_key|
  command_yml = commands_configs[command_config_key]
  if command_yml['command_line'].nil?
    $logger.warn "Command #{command_config_key} does not contain a " +
                 "command_line spec - ignoring"
    next
  end
  if command_yml['cron'].nil?
    $logger.warn "Command #{command_config_key} does not contain a " +
                 "cron spec - ignoring"
    next
  end
  _timeout = command_yml.key?('timeout') ? command_yml['timeout'] : DEFAULT_TIMEOUT
  commands.push(Command.new(command_config_key,
                            command_yml['command_line'], command_yml['cron'], _timeout))
end

# Log commands we're going to run.
$logger.info "Running on: #{commands.inspect}"

# Rotate on all commands
loop do
  timestamp = Time.now.to_i
  commands.each do |command|
    next unless (timestamp % command.execution_time_in_sec) == 0
    pid = Process.fork do
      command.run
    end
    Process.detach(pid)
  end
  sleep 1
end

# rubocop:enable all
