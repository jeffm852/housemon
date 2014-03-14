stream = require 'stream'
fs = require 'fs'

LOGGER_PATH = './logger-old'

mkdir = (path) ->
  unless fs.existsSync path
    fs.mkdirSync path

dateFilename = (now) ->
  # construct the date value as 8 digits
  y = now.getUTCFullYear()
  d = (y * 100 + now.getUTCMonth() + 1) * 100 + now.getUTCDate()
  # then massage it as a string to produce a file name
  path = "#{LOGGER_PATH}/#{y}"
  mkdir path
  path + "/#{d}.txt"

timeString = (now) ->
  date = new Date(now)
  # first construct the value as 10 digits
  digits = date.getUTCMilliseconds() + 1000 *
          (date.getUTCSeconds() + 100 *
          (date.getUTCMinutes() + 100 *
          (date.getUTCHours() + 100)))
  # then massage it as a string to get the punctuation right
  digits.toString().replace /.(..)(..)(..)(...)/, '$1:$2:$3.$4'

class Logger extends stream.Writable
  constructor: ->
    super objectMode: true
    mkdir LOGGER_PATH
    @currDate = null

  _write: (data, encoding, done) ->
    if data?
      now = new Date
      unless now.getUTCDate() is @currDate
        @currDate = now.getUTCDate()
        fs.closeSync @fd  if @fd?
        logname = dateFilename now
        console.info 'Logger: writing to', logname
        @fd = fs.openSync logname, 'a'
      # L 01:02:03.537 usb-A40117UK OK 9 25 54 66 235 61 210 226 33 19
      {time,dev,msg} = data
      msg = "L #{timeString time ? now} #{dev} #{msg}\n"
      fs.write @fd, Buffer(msg), 0, msg.length, null, done
    else
      fs.closeSync @fd  if @fd?
      done()

module.exports = (app, plugin) ->
  app.register 'sink.logger', Logger
