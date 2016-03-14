'use strict';

var botUtilities = require('bot-utilities');
var execFile = require('child_process').execFile;
var program = require('commander');
var Twit = require('twit');
var _ = require('lodash');

var Canvas = require('canvas-utilities').Canvas;
var canvasUtilities = require('canvas-utilities/lib/utilities.js');
var ImageHelper = require('canvas-utilities/lib/image-helper.js');

_.mixin(botUtilities.lodashMixins);
_.mixin(Twit.prototype, botUtilities.twitMixins);

function makeImage(cb) {
  var WIDTH = 1024;
  var HEIGHT = 1024;

  var canvas = new Canvas(WIDTH, HEIGHT);
  var ctx = canvasUtilities.getContext(canvas);

  ImageHelper.fromFile('./out.png').context(ctx).draw(0, 0, WIDTH, HEIGHT);

  canvas.toBuffer(function (err, buffer) {
    if (err) {
      throw err;
    }

    cb(err, buffer);
  });
}

program
  .command('tweet')
  .description('Generate and tweet an image')
  .option('-r, --random', 'only post a percentage of the time')
  .action(botUtilities.randomCommand(function () {
    execFile('./render', [], (err, stdout) => {
      if (err) {
        throw err;
      }

      console.log(stdout);

      makeImage(function (err, buffer) {
        var T = new Twit(botUtilities.getTwitterAuthFromEnv());

        var tweet = '';

        if (_.percentChance(25)) {
          var bot = botUtilities.imageBot();

          tweet += botUtilities.heyYou(bot);

          if (bot === 'Lowpolybot') {
            tweet += ' #noRT';
          }
        }

        tweet = {status: tweet};

        T.updateWithMedia(tweet, buffer, function (err, response) {
          if (err) {
            return console.error('TUWM error', err, response.statusCode);
          }

          console.log('TUWM OK');
        });
      });
    });
  }));

program.parse(process.argv);
