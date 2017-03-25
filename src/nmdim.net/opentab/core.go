//core.go
//
//Copyright 2017-present Zhang Yuxing. All Rights Reserved.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package opentab

import (
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// ssl证书 文件存储= master_api  owner=  _File= acl过滤=  严格模式  scan  类封装 批处理  time= 字段过滤=   ip黑名单防恶意刷新   日志  性能记录  数据迁移  备份  conut=  更新缓存

func Run() {
	r := gin.Default()
	loadMiddleware(r)
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type","text/html;charset=utf-8")
		c.String(http.StatusOK, `<script>/**
 * segment - A little JavaScript class (without dependencies) to draw and animate SVG path strokes
 * @version v1.0.8
 * @link https://github.com/lmgonzalves/segment
 * @license MIT
 */
function Segment(i,e,t,n){this.path=i,this.length=i.getTotalLength(),this.path.style.strokeDashoffset=2*this.length,this.begin="undefined"!=typeof e?this.valueOf(e):0,this.end="undefined"!=typeof t?this.valueOf(t):this.length,this.circular="undefined"!==n?n:!1,this.timer=null,this.animationTimer=null,this.draw(this.begin,this.end,0,{circular:this.circular})}!function(){for(var i=0,e=["ms","moz","webkit","o"],t=0;t<e.length&&!window.requestAnimationFrame;++t)window.requestAnimationFrame=window[e[t]+"RequestAnimationFrame"],window.cancelAnimationFrame=window[e[t]+"CancelAnimationFrame"]||window[e[t]+"CancelRequestAnimationFrame"];window.requestAnimationFrame||(window.requestAnimationFrame=function(e){var t=(new Date).getTime(),n=Math.max(0,16-(t-i)),s=window.setTimeout(function(){e(t+n)},n);return i=t+n,s}),window.cancelAnimationFrame||(window.cancelAnimationFrame=function(i){clearTimeout(i)})}(),Segment.prototype={draw:function(i,e,t,n){if(this.circular=n&&n.hasOwnProperty("circular")?n.circular:!1,t){var s=n&&n.hasOwnProperty("delay")?1e3*parseFloat(n.delay):0,h=n&&n.hasOwnProperty("easing")?n.easing:null,r=n&&n.hasOwnProperty("callback")?n.callback:null,a=this;if(this.stop(),s)return delete n.delay,this.timer=setTimeout(function(){a.draw(i,e,t,n)},s),this.timer;var l=new Date,o=this.begin,g=this.end,c=this.valueOf(i),u=this.valueOf(e);!function d(){var i=new Date,e=(i-l)/1e3,n=e/parseFloat(t),s=n;return"function"==typeof h&&(s=h(s)),n>1?s=1:a.animationTimer=window.requestAnimationFrame(d),a.begin=o+(c-o)*s,a.end=g+(u-g)*s,a.begin=a.begin<0&&!a.circular?0:a.begin,a.begin=a.begin>a.length&&!a.circular?a.length:a.begin,a.end=a.end<0&&!a.circular?0:a.end,a.end=a.end>a.length&&!a.circular?a.length:a.end,a.end-a.begin<=a.length&&a.end-a.begin>0?a.draw(a.begin,a.end,0,{circular:a.circular}):a.circular&&a.end-a.begin>a.length?a.draw(0,a.length,0,{circular:a.circular}):a.draw(a.begin+(a.end-a.begin),a.end-(a.end-a.begin),0,{circular:a.circular}),n>1&&"function"==typeof r?r.call(a):void 0}()}else this.path.style.strokeDasharray=this.strokeDasharray(i,e)},strokeDasharray:function(i,e){if(this.begin=this.valueOf(i),this.end=this.valueOf(e),this.circular){var t=this.begin>this.end||this.begin<0&&this.begin<-1*this.length?parseInt(this.begin/parseInt(this.length)):parseInt(this.end/parseInt(this.length));0!==t&&(this.begin=this.begin-this.length*t,this.end=this.end-this.length*t)}if(this.end>this.length){var n=this.end-this.length;return[this.length,this.length,n,this.begin-n,this.end-this.begin].join(" ")}if(this.begin<0){var s=this.length+this.begin;return this.end<0?[this.length,this.length+this.begin,this.end-this.begin,s-this.end,this.end-this.begin,this.length].join(" "):[this.length,this.length+this.begin,this.end-this.begin,s-this.end,this.length].join(" ")}return[this.length,this.length+this.begin,this.end-this.begin].join(" ")},valueOf:function(i){var e=parseFloat(i);if(("string"==typeof i||i instanceof String)&&~i.indexOf("%")){var t;~i.indexOf("+")?(t=i.split("+"),e=this.percent(t[0])+parseFloat(t[1])):~i.indexOf("-")?(t=i.split("-"),e=3===t.length?-this.percent(t[1])-parseFloat(t[2]):t[0]?this.percent(t[0])-parseFloat(t[1]):-this.percent(t[1])):e=this.percent(i)}return e},stop:function(){window.cancelAnimationFrame(this.animationTimer),this.animationTimer=null,clearTimeout(this.timer),this.timer=null},percent:function(i){return parseFloat(i)/100*this.length}};</script>
		<script>(function (global, factory) {
  typeof exports === 'object' && typeof module !== 'undefined' ? factory(exports) :
  typeof define === 'function' && define.amd ? define(['exports'], factory) :
  (factory((global.d3_ease = {})));
}(this, function (exports) { 'use strict';

  var linearIn = {
    ease: function(t) {
      return +t;
    }
  };

  var quadIn = {
    ease: function(t) {
      return t * t;
    }
  };

  var quadOut = {
    ease: function(t) {
      return t * (2 - t);
    }
  };

  var quadInOut = {
    ease: function(t) {
      return ((t *= 2) <= 1 ? t * t : --t * (2 - t) + 1) / 2;
    }
  };

  var cubicIn = {
    ease: function(t) {
      return t * t * t;
    }
  };

  var cubicOut = {
    ease: function(t) {
      return --t * t * t + 1;
    }
  };

  var cubicInOut = {
    ease: function(t) {
      return ((t *= 2) <= 1 ? t * t * t : (t -= 2) * t * t + 2) / 2;
    }
  };

  var exponent = 3;

  var polyIn = (function polyIn(e) {
    return e = +e, {
      exponent: polyIn,
      ease: function(t) {
        return Math.pow(t, e);
      }
    };
  })(exponent);

  var polyOut = (function polyOut(e) {
    return e = +e, {
      exponent: polyOut,
      ease: function(t) {
        return 1 - Math.pow(1 - t, e);
      }
    };
  })(exponent);

  var polyInOut = (function polyInOut(e) {
    return e = +e, {
      exponent: polyInOut,
      ease: function(t) {
        return ((t *= 2) <= 1 ? Math.pow(t, e) : 2 - Math.pow(2 - t, e)) / 2;
      }
    };
  })(exponent);

  var pi = Math.PI;
  var halfPi = pi / 2;
  var sinIn = {
    ease: function(t) {
      return 1 - Math.cos(t * halfPi);
    }
  };

  var sinOut = {
    ease: function(t) {
      return Math.sin(t * halfPi);
    }
  };

  var sinInOut = {
    ease: function(t) {
      return (1 - Math.cos(pi * t)) / 2;
    }
  };

  var expIn = {
    ease: function(t) {
      return Math.pow(2, 10 * t - 10);
    }
  };

  var expOut = {
    ease: function(t) {
      return 1 - Math.pow(2, -10 * t);
    }
  };

  var expInOut = {
    ease: function(t) {
      return ((t *= 2) <= 1 ? Math.pow(2, 10 * t - 10) : 2 - Math.pow(2, 10 - 10 * t)) / 2;
    }
  };

  var circleIn = {
    ease: function(t) {
      return 1 - Math.sqrt(1 - t * t);
    }
  };

  var circleOut = {
    ease: function(t) {
      return Math.sqrt(1 - --t * t);
    }
  };

  var circleInOut = {
    ease: function(t) {
      return ((t *= 2) <= 1 ? 1 - Math.sqrt(1 - t * t) : Math.sqrt(1 - (t -= 2) * t) + 1) / 2;
    }
  };

  var b1 = 4 / 11;
  var b2 = 6 / 11;
  var b3 = 8 / 11;
  var b4 = 3 / 4;
  var b5 = 9 / 11;
  var b6 = 10 / 11;
  var b7 = 15 / 16;
  var b8 = 21 / 22;
  var b9 = 63 / 64;
  var b0 = 1 / b1 / b1;
  function bounce(t) {
    return (t = +t) < b1 ? b0 * t * t : t < b3 ? b0 * (t -= b2) * t + b4 : t < b6 ? b0 * (t -= b5) * t + b7 : b0 * (t -= b8) * t + b9;
  }

  var bounceIn = {
    ease: function(t) {
      return 1 - bounce(1 - t);
    }
  };

  var bounceOut = {
    ease: bounce
  };

  var bounceInOut = {
    ease: function(t) {
      return ((t *= 2) <= 1 ? 1 - bounce(1 - t) : bounce(t - 1) + 1) / 2;
    }
  };

  var overshoot = 1.70158;

  var backIn = (function backIn(s) {
    return s = +s, {
      overshoot: backIn,
      ease: function(t) {
        return t * t * ((s + 1) * t - s);
      }
    };
  })(overshoot);

  var backOut = (function backOut(s) {
    return s = +s, {
      overshoot: backOut,
      ease: function(t) {
        return --t * t * ((s + 1) * t + s) + 1;
      }
    };
  })(overshoot);

  var backInOut = (function backInOut(s) {
    return s = +s, {
      overshoot: backInOut,
      ease: function(t) {
        return ((t *= 2) < 1 ? t * t * ((s + 1) * t - s) : (t -= 2) * t * ((s + 1) * t + s) + 2) / 2;
      }
    };
  })(overshoot);

  var tau = 2 * Math.PI;
  var amplitude = 1;
  var period = 0.3;
  var elasticIn = (function elasticIn(a, p) {
    var s = Math.asin(1 / (a = Math.max(1, a))) * (p /= tau);
    return {
      amplitude: function(a) { return elasticIn(a, p * tau); },
      period: function(p) { return elasticIn(a, p); },
      ease: function(t) {
        return a * Math.pow(2, 10 * --t) * Math.sin((s - t) / p);
      }
    };
  })(amplitude, period);

  var elasticOut = (function elasticOut(a, p) {
    var s = Math.asin(1 / (a = Math.max(1, a))) * (p /= tau);
    return {
      amplitude: function(a) { return elasticOut(a, p * tau); },
      period: function(p) { return elasticOut(a, p); },
      ease: function(t) {
        return 1 - a * Math.pow(2, -10 * (t = +t)) * Math.sin((t + s) / p);
      }
    };
  })(amplitude, period);

  var elasticInOut = (function elasticInOut(a, p) {
    var s = Math.asin(1 / (a = Math.max(1, a))) * (p /= tau);
    return {
      amplitude: function(a) { return elasticInOut(a, p * tau); },
      period: function(p) { return elasticInOut(a, p); },
      ease: function(t) {
        return ((t = t * 2 - 1) < 0
            ? a * Math.pow(2, 10 * t) * Math.sin((s - t) / p)
            : 2 - a * Math.pow(2, -10 * t) * Math.sin((s + t) / p)) / 2;
      }
    };
  })(amplitude, period);

  var version = "0.6.0";

  exports.version = version;
  exports.easeLinearIn = linearIn;
  exports.easeLinearOut = linearIn;
  exports.easeLinearInOut = linearIn;
  exports.easeQuadIn = quadIn;
  exports.easeQuadOut = quadOut;
  exports.easeQuadInOut = quadInOut;
  exports.easeCubicIn = cubicIn;
  exports.easeCubicOut = cubicOut;
  exports.easeCubicInOut = cubicInOut;
  exports.easePolyIn = polyIn;
  exports.easePolyOut = polyOut;
  exports.easePolyInOut = polyInOut;
  exports.easeSinIn = sinIn;
  exports.easeSinOut = sinOut;
  exports.easeSinInOut = sinInOut;
  exports.easeExpIn = expIn;
  exports.easeExpOut = expOut;
  exports.easeExpInOut = expInOut;
  exports.easeCircleIn = circleIn;
  exports.easeCircleOut = circleOut;
  exports.easeCircleInOut = circleInOut;
  exports.easeBounceIn = bounceIn;
  exports.easeBounceOut = bounceOut;
  exports.easeBounceInOut = bounceInOut;
  exports.easeBackIn = backIn;
  exports.easeBackOut = backOut;
  exports.easeBackInOut = backInOut;
  exports.easeElasticIn = elasticIn;
  exports.easeElasticOut = elasticOut;
  exports.easeElasticInOut = elasticInOut;

}));</script>
		<script>function Letters(el, options){
    this.el = el;
    this.text = el.innerHTML.toLowerCase();
    el.innerHTML = '<span style="position: absolute; visibility: hidden;">'+ this.text +'</span>';
    this.init(options);
}

Letters.prototype = {
    init: function(options){
        this.size = options && options.hasOwnProperty('size') ? options.size : 100;
        this.weight = options && options.hasOwnProperty('weight') ? options.weight : 1;
        this.rounded = options && options.hasOwnProperty('rounded') ? options.rounded : false;
        this.color = options && options.hasOwnProperty('color') ? options.color : '#5F6062';
        this.duration = options && options.hasOwnProperty('duration') ? options.duration : 1;
        this.delay = options && options.hasOwnProperty('delay') ? options.delay : [0, 0.05];
        this.fade = options && options.hasOwnProperty('fade') ? options.fade : 0.5;
        this.easing = options && options.hasOwnProperty('easing') ? options.easing : d3_ease.easeCubicInOut.ease;
        this.individualDelays = options && options.hasOwnProperty('individualDelays') ? options.individualDelays : false;
        this.visible = false;
        this._svgNS = 'http://www.w3.org/2000/svg';
        this.svgs = [];
        this.timer = [];
        this.segments = [];
        this.drawText();
    },

    drawText: function(){
        this.letters = this.text.replace(/[^a-z ]/g, '').split('');
        var self = this;
        var size = typeof self.size === 'number' ? self.size : self.size[0];
        var weight = typeof self.weight === 'number' ? self.weight : self.weight[0];
        var rounded = typeof self.rounded === 'boolean' ? self.rounded : self.rounded[0];
        var color = typeof self.color === 'string' ? self.color : self.color[0];
        var fade = typeof self.fade === 'number' ? self.fade : self.fade[0];
        self.renderAlphabet = typeof self.weight === 'number' && typeof self.rounded === 'boolean' ? self.initAlphabet(weight, rounded) && false : true;
        var spaces = 0;
        this.letters.forEach(function(element, index){
            size = typeof self.size[index] === 'number' ? self.size[index] : size;
            weight = typeof self.weight[index] === 'number' ? self.weight[index] : weight;
            rounded = typeof self.rounded[index] === 'boolean' ? self.rounded[index] : rounded;
            color = typeof self.color !== 'string' && typeof self.color[index] === 'string' ? self.color[index] : color;
            fade = typeof self.fade[index] === 'number' ? self.fade[index] : fade;
            if(self.renderAlphabet){
                self.initAlphabet(weight, rounded);
            }
            if(element === ' '){
                self.drawLetter(element, index, size);
                spaces++;
            }else{
                self.segments[index - spaces] = self.drawLetter(element, index - spaces, size, weight, rounded, color, fade);
            }
        });
        this.letters = this.text.replace(/[^a-z]/g, '').split('');
    },

    drawLetter: function(letter, index, size, weight, rounded, color, fade){
        if(letter === ' '){
            this.createSVG(this.alphabet[letter].svg, 'space', index, size);
        }else{
            var svg = this.createSVG(this.alphabet[letter].svg, letter, index, size, fade);
            var segments = [];
            var pathsD = this.alphabet[letter].paths;
            var initial = this.alphabet[letter].initial;
            var path;
            var segment;
            var self = this;
            pathsD.forEach(function(element, index){
                path = document.createElementNS(self._svgNS, "path");
                path.setAttribute('d', element);
                path.setAttribute('stroke-width', weight);
                path.setAttribute('stroke-linecap', rounded ? 'round' : 'butt');
                path.setAttribute('stroke', color);
                path.setAttribute('fill', 'none');
                segment = new Segment(path, initial[index].begin, initial[index].end, true);
                segments.push(segment);
                svg.appendChild(path);
            });
            return segments;
        }
    },

    createSVG: function(svgSize, letter, index, size, fade){
        var svg = document.createElementNS(this._svgNS, "svg");
        svg.setAttribute('aria-hidden', 'true');
        svg.setAttribute('role', 'img');
        svg.setAttribute('viewBox', '0 0 '+ svgSize.width +' '+ svgSize.height);
        svg.setAttribute('height', size + 'px');
        var width = size * (svgSize.width / svgSize.height);
        svg.setAttribute('width', Math.ceil(width) + 'px');
        svg.setAttribute('class', 'letter letter-' + letter + (letter !== 'space' ? ' letter-' + (index + 1) : ''));
        this.el.appendChild(svg);
        if(letter !== 'space'){
            this.svgs[index] = svg;
            if(fade){
                svg.style.opacity = 0;
                this.setupFade(fade, index);
            }
            return svg;
        }
    },

    setupFade: function(fade, index){
        var svg = this.svgs[index];
        if(fade){
            svg.style.transition = fade + 's opacity';
            svg.style.webkitTransition = fade + 's opacity';
        }else{
            svg.style.transition = null;
        }
    },

    show: function(options){
        this.visible = true;
        this.clear();
        var _duration = options && options.hasOwnProperty('duration') ? options.duration : this.duration;
        var _delay = options && options.hasOwnProperty('delay') ? options.delay : this.delay;
        var _fade = options && options.hasOwnProperty('fade') ? options.fade : this.fade;
        var _easing = options && options.hasOwnProperty('easing') ? options.easing : this.easing;
        var _individualDelays = options && options.hasOwnProperty('individualDelays') ? options.individualDelays : this.individualDelays;
        var tempDelay = 0;
        var lastDelay = 0;
        var lastDuration = typeof _duration === 'number' ? _duration : _duration[0];
        var lastFade = typeof _fade === 'number' ? _fade : _fade[0];
        var lastEasing = typeof _easing === 'function' ? _easing : _easing[0];
        var self = this;
        var weight, rounded;
        if(self.renderAlphabet){
            weight = typeof self.weight === 'number' ? self.weight : self.weight[0];
            rounded = typeof self.rounded === 'boolean' ? self.rounded : self.rounded[0];
        }
        var length = this.letters.length;
        this.letters.forEach(function(element, index){
            lastDuration = typeof _duration[index] === 'number' ? _duration[index] : lastDuration;
            lastFade = typeof _fade[index] === 'number' ? _fade[index] : lastFade;
            lastEasing = typeof _easing[index] === 'function' ? _easing[index] : lastEasing;
            if(typeof _delay === 'number'){
                tempDelay = _individualDelays ? _delay : tempDelay + _delay;
            }else{
                lastDelay = typeof _delay[index] === 'number' ? _delay[index] : lastDelay;
                tempDelay = _individualDelays ? lastDelay : tempDelay + lastDelay;
            }
            if(tempDelay){
                self.timer[index] = setTimeout(function(){
                    self.setupFade(lastFade, index);
                    self.svgs[index].style.opacity = 1;
                    self.el.offsetHeight;
                }, 1000 * tempDelay);
            }else{
                self.setupFade(lastFade, index);
                self.svgs[index].style.opacity = 1;
                self.el.offsetHeight;
            }
            if(self.renderAlphabet){
                weight = typeof self.weight[index] === 'number' ? self.weight[index] : weight;
                rounded = typeof self.rounded[index] === 'boolean' ? self.rounded[index] : rounded;
                self.initAlphabet(weight, rounded);
            }
            var _callback = index === length - 1 && options && options.hasOwnProperty('callback') ? options.callback : null;
            var finals = self.alphabet[element].final;
            var segments = self.segments[index];
            segments.forEach(function(element, index){
                element.draw(finals[index].begin, finals[index].end, lastDuration, {delay: tempDelay, circular: true, easing: lastEasing, callback: _callback});
                _callback = null;
            });
        });
    },

    hide: function(options){
        this.visible = false;
        this.clear();
        var _duration = options && options.hasOwnProperty('duration') ? options.duration : this.duration;
        var _delay = options && options.hasOwnProperty('delay') ? options.delay : this.delay;
        var _fade = options && options.hasOwnProperty('fade') ? options.fade : this.fade;
        var _easing = options && options.hasOwnProperty('easing') ? options.easing : this.easing;
        var _individualDelays = options && options.hasOwnProperty('individualDelays') ? options.individualDelays : this.individualDelays;
        var tempDelay = 0;
        var lastDelay = 0;
        var lastDuration = typeof _duration === 'number' ? _duration : _duration[0];
        var lastFade = typeof _fade === 'number' ? _fade : _fade[0];
        var lastEasing = typeof _easing === 'function' ? _easing : _easing[0];
        var self = this;
        var length = this.letters.length;
        this.letters.forEach(function(element, index){
            lastDuration = typeof _duration[index] === 'number' ? _duration[index] : lastDuration;
            lastFade = typeof _fade[index] === 'number' ? _fade[index] : lastFade;
            lastEasing = typeof _easing[index] === 'function' ? _easing[index] : lastEasing;
            if(typeof _delay === 'number'){
                tempDelay = _individualDelays ? _delay : tempDelay + _delay;
            }else{
                lastDelay = typeof _delay[index] === 'number' ? _delay[index] : lastDelay;
                tempDelay = _individualDelays ? lastDelay : tempDelay + lastDelay;
            }
            self.setupFade(lastFade, index);
            var time = 1000 * (lastDuration - lastFade + tempDelay);
            if(time){
                self.timer[index] = setTimeout(function(){
                    self.svgs[index].style.opacity = 0;
                    self.el.offsetHeight;
                }, time);
            }else{
                self.svgs[index].style.opacity = 0;
                self.el.offsetHeight;
            }
            var _callback = index === length - 1 && options && options.hasOwnProperty('callback') ? options.callback : null;
            var initial = self.alphabet[element].initial;
            var segments = self.segments[index];
            segments.forEach(function(element, index){
                element.draw(initial[index].begin, initial[index].end, lastDuration, {delay: tempDelay, circular: true, easing: lastEasing, callback: _callback});
                _callback = null;
            });
        });
    },

    toggle: function(options){
        if(this.visible){
            this.hide(options);
        }else{
            this.show(options);
        }
    },

    showInstantly: function(){
        this.show({duration: 0, delay: 0, fade: 0});
    },

    hideInstantly: function(){
        this.hide({duration: 0, delay: 0, fade: 0});
    },

    toggleInstantly: function(){
        this.toggle({duration: 0, delay: 0, fade: 0});
    },

    clear: function(){
        var self = this;
        this.letters.forEach(function(element, index){
            clearTimeout(self.timer[index]);
        });
    },

    initAlphabet: function(weight, rounded){
        var width = 44 + weight;
        var height = 94 + (2 * weight);
        var radius = 20;
        var _weight = rounded ? 0 : weight / 2;
        var circle = 'm 0 -'+ radius +' a '+ radius +' '+ radius +' 0 1 1 0 '+ (2 * radius) +' a '+ radius +' '+ radius +' 0 1 1 0 -'+ (2 * radius);
        var circleCenter = 'm '+ (width / 2) +' '+ (height / 2) +' ' + circle;
        var circleLeft = 'm '+ (width / 2) +' '+ (height / 2) +' ' + circle;
        var circleRight = 'm '+ ((width / 2) + (2 * radius)) +' '+ (height / 2) +' ' + circle;
        var circleF = 'm '+ (width / 2) +' '+ ((height / 2) - 12) +' ' + circle;
        var circleF2 = 'm '+ (width / 2) +' '+ ((height / 2) + (weight / 2)) +' ' + circle;
        var circleGY = 'm '+ (width / 2) +' '+ ((height / 2) + radius + _weight) +' ' + circle;
        var circleJ = 'm '+ ((width / 2) - radius - _weight) +' '+ ((height / 2) + radius + _weight) +' ' + circle;
        var circleSXZ = 'm '+ ((width / 2) - radius - _weight) +' '+ (height / 2) +' ' + circle;
        var circleSXZ2 = 'm '+ ((width / 2) + radius - _weight) +' '+ (height / 2) +' ' + circle;
        var lineLeft = 'm '+ ((width / 2) - radius) +' 0 l 0 '+ height;
        var lineLeftLarge = 'm '+ ((width / 2) - radius) +' 0 l 0 '+ height;
        var lineRight = 'm '+ ((width / 2) + radius) +' 0 l 0 '+ height;
        var lineRightLarge = 'm '+ ((width / 2) + (3 * radius)) +' 0 l 0 '+ height;
        var lineCenter = 'm '+ ((width / 2) - radius) +' 0 l 0 '+ height;
        var lineE = 'm 0 '+ ((height / 2) + _weight + (_weight ? 0 : 1)) +' l '+ width +' 0';
        var lineJ = 'm '+ ((width / 2) - _weight) +' 0 l 0 '+ height;
        var lineT = 'm 0 '+ ((height / 2) - radius) +' l '+ width +' 0';
        var svgSmall = {width: width, height: height};
        var svgLarge = {width: (width + (2 * radius)), height: height};
        var svgFJRT = {width: (width - radius - _weight), height: height};
        var svgSXZ = {width: (width - (2 * _weight)), height: height};
        var p50plus23 = '50% + ' + (radius + _weight);
        var p50minus23 = '50% - ' + (radius + _weight);
        var p50minus20 = '50% - ' + radius;
        var p50minus35 = '50% - ' + (radius + _weight + 12);
        this.alphabet = {
            'a': {
                svg: svgSmall,
                paths: [circleCenter, lineRight],
                initial: [{begin: '5%', end: '5%'}, {begin: '90%', end: '90%'}],
                final: [{begin: '50%', end: '125% + 1'}, {begin: '50%', end: p50plus23}]
            },
            'b': {
                svg: svgSmall,
                paths: [circleCenter, lineLeft],
                initial: [{begin: '-125%', end: '-125%'}, {begin: '5%', end: '5%'}],
                final: [{begin: '-125%', end: '-50%'}, {begin: p50minus35, end: p50plus23}]
            },
            'c': {
                svg: svgSmall,
                paths: [circleCenter],
                initial: [{begin: '105%', end: '105%'}],
                final: [{begin: '-50%', end: '25%'}]
            },
            'd': {
                svg: svgSmall,
                paths: [circleCenter, lineRight],
                initial: [{begin: '125%', end: '125%'}, {begin: '5%', end: '5%'}],
                final: [{begin: '50%', end: '125%'}, {begin: p50minus35, end: p50plus23}]
            },
            'e': {
                svg: svgSmall,
                paths: [circleCenter, lineE],
                initial: [{begin: '-105%', end: '-105%'}, {begin: '5%', end: '5%'}],
                final: [{begin: '50%', end: '125% + 1'}, {begin: '50%', end: p50plus23}]
            },
            'f': {
                svg: svgFJRT,
                paths: [circleF, circleF2, lineLeft],
                initial: [{begin: '-25%', end: '-25%'}, {begin: '1', end: '1'}, {begin: '50% - 13', end: '50% - 13'}],
                final: [{begin: '-25%', end: '1'}, {begin: '-25%', end: '1'}, {begin: '50% - 13', end: p50plus23}]
            },
            'g': {
                svg: svgSmall,
                paths: [circleCenter, circleGY, lineRight],
                initial: [{begin: '25%', end: '25%'}, {begin: '-10%', end: '-10%'}, {begin: '50%', end: '50%'}],
                final: [{begin: '-50%', end: '25% + 1'}, {begin: '-75% - 1', end: '-50%'}, {begin: '50%', end: p50plus23}]
            },
            'h': {
                svg: svgSmall,
                paths: [circleCenter, lineLeft, lineRight],
                initial: [{begin: '-25%', end: '-25%'}, {begin: '10%', end: '10%'}, {begin: '90%', end: '90%'}],
                final: [{begin: '-150%', end: '-75% + 1'}, {begin: p50minus35, end: '50%'}, {begin: '50%', end: p50plus23}]
            },
            'i': {
                svg: {width: (width - (2 * radius)), height: height},
                paths: [lineCenter],
                initial: [{begin: '50%', end: '50%'}],
                final: [{begin: p50minus23, end: p50plus23}]
            },
            'j': {
                svg: svgFJRT,
                paths: [circleJ, lineJ],
                initial: [{begin: '5%', end: '5%'}, {begin: p50minus23, end: p50minus23}],
                final: [{begin: '25% - 1', end: '50%'}, {begin: p50minus23, end: p50plus23}]
            },
            'k': {
                svg: svgSmall,
                paths: [circleCenter, lineLeft],
                initial: [{begin: '50%', end: '50%'}, {begin: '80%', end: '80%'}],
                final: [{begin: '25%', end: '100%'}, {begin: p50minus35, end: p50plus23}]
            },
            'l': {
                svg: {width: (width - (2 * radius)), height: height},
                paths: [lineCenter],
                initial: [{begin: '50%', end: '50%'}],
                final: [{begin: p50minus35, end: p50plus23}]
            },
            'm': {
                svg: svgLarge,
                paths: [circleLeft, circleRight, lineLeftLarge],
                initial: [{begin: '90%', end: '90%'}, {begin: '90%', end: '90%'}, {begin: '90%', end: '90%'}],
                final: [{begin: '-25%', end: '50%'}, {begin: '-25%', end: '50%'}, {begin: p50minus23, end: p50plus23}]
            },
            'n': {
                svg: svgSmall,
                paths: [circleCenter, lineLeft],
                initial: [{begin: '-60%', end: '-60%'}, {begin: '90%', end: '90%'}],
                final: [{begin: '-25%', end: '50%'}, {begin: p50minus23, end: p50plus23}]
            },
            'o': {
                svg: svgSmall,
                paths: [circleCenter],
                initial: [{begin: '-150%', end: '-150%'}],
                final: [{begin: '2%', end: '98%'}]
            },
            'p': {
                svg: svgSmall,
                paths: [circleCenter, lineLeft],
                initial: [{begin: '25%', end: '25%'}, {begin: p50plus23, end: p50plus23}],
                final: [{begin: '-25% - 1', end: '50%'}, {begin: '50%', end: '50% + 46'}]
            },
            'q': {
                svg: svgSmall,
                paths: [circleCenter, lineRight],
                initial: [{begin: '-25%', end: '-25%'}, {begin: p50plus23, end: p50plus23}],
                final: [{begin: '-50%', end: '25% + 1'}, {begin: '50%', end: '50% + 46'}]
            },
            'r': {
                svg: svgFJRT,
                paths: [circleCenter, lineLeft],
                initial: [{begin: '-25% - 1', end: '-25% - 1'}, {begin: '50%', end: '50%'}],
                final: [{begin: '-25% - 1', end: '1'}, {begin: '50%', end: p50plus23}]
            },
            's': {
                svg: svgSXZ,
                paths: [circleSXZ, circleSXZ2],
                initial: [{begin: '50%', end: '50%'}, {begin: '1', end: '1'}],
                final: [{begin: '25%', end: '50%'}, {begin: '-25% - 1', end: '1'}]
            },
            't': {
                svg: svgFJRT,
                paths: [circleCenter, lineLeft, lineT],
                initial: [{begin: '50%', end: '50%'}, {begin: '5%', end: '5%'}, {begin: '50%', end: '50%'}],
                final: [{begin: '50%', end: '75% + 1'}, {begin: p50minus35, end: '50%'}, {begin: p50minus20, end: '50%'}]
            },
            'u': {
                svg: svgSmall,
                paths: [circleCenter, lineRight],
                initial: [{begin: '50%', end: '50%'}, {begin: '80%', end: '80%'}],
                final: [{begin: '25%', end: '100%'}, {begin: p50minus23, end: p50plus23}]
            },
            'v': {
                svg: svgSmall,
                paths: [circleCenter, lineRight],
                initial: [{begin: '10%', end: '10%'}, {begin: '20%', end: '20%'}],
                final: [{begin: '25% - 1', end: '100%'}, {begin: p50minus23, end: '50%'}]
            },
            'w': {
                svg: svgLarge,
                paths: [circleLeft, circleRight, lineRightLarge],
                initial: [{begin: '10%', end: '10%'}, {begin: '10%', end: '10%'}, {begin: '20%', end: '20%'}],
                final: [{begin: '25%', end: '100%'}, {begin: '25% - 1', end: '100%'}, {begin: p50minus23, end: '50%'}]
            },
            'x': {
                svg: svgSXZ,
                paths: [circleSXZ, circleSXZ2],
                initial: [{begin: '50%', end: '50%'}, {begin: '1', end: '1'}],
                final: [{begin: '0', end: '50%'}, {begin: '-50%', end: '1'}]
            },
            'y': {
                svg: svgSmall,
                paths: [circleCenter, circleGY, lineRight],
                initial: [{begin: '1', end: '1'}, {begin: '1', end: '1'}, {begin: p50minus23, end: p50minus23}],
                final: [{begin: '-75%', end: '1'}, {begin: '25% - 1', end: '50%'}, {begin: p50minus23, end: p50plus23}]
            },
            'z': {
                svg: svgSXZ,
                paths: [circleSXZ, circleSXZ2],
                initial: [{begin: '25%', end: '25%'}, {begin: '75%', end: '75%'}],
                final: [{begin: '0', end: '25% + 1'}, {begin: '50%', end: '75%'}]
            },
            ' ': {
                svg: svgSmall
            }
        };
    }
};</script>
<div class="c">
	<div class="text">Opentab</div>
</div>
		<script>
		//选择元素
var el = document.querySelector('.text');

//每个选项可以定义为单个也可以定义为数组
var options = {
    size: 200,         // 字体大小，决定文字的高度，px
    weight: 5,         // 粗体，px
    rounded: false,    // 字母结尾圆角
    color: ['#f90','#cccccc'],  // 字体颜色
    duration: 1,       // 每个字母的动画展示时长 (seconds)
    delay: [0, 0.1],  // 每个字母间的动画延时
    fade: 0.5,         // 渐显效果的时长(seconds)
    easing: d3_ease.easeCubicInOut.ease,   // 缓冲动画效果
    individualDelays: false     // 默认false，如果设置成false，则动画效果会从左到右过度展示，如果是true，字母动画效果是同步展示
};

var myText = new Letters(el, options);

myText.show();
		</script>
<style>
	.c{
		justify-content: center;
		align-items: center;
		max-width: calc(100vw);
    		height: calc(100vh);
		display: flex;
		background-color:#2C3437;
	}
</style>`)
	})
	loadRouter(r)
	fmt.Println(`应用已启动`)

	appInfo:=GetAppInfo()

	fmt.Println(`AppId:`+appInfo[`AppId`])
	fmt.Println(`AppKey:`+appInfo[`AppKey`])
	fmt.Println(`MasterKey:`+appInfo[`MasterKey`])

	if App.IsSSL {
		r.RunTLS(`:`+App.ListenPosrt, If(App.IsLinux, "/usr/bin/", "").(string)+App.CertFile, If(App.IsLinux, "/usr/bin/", "").(string)+App.KeyFile)
	} else {
		r.Run(`:` + App.ListenPosrt)
	}

}

func loadMiddleware(r *gin.Engine)  {
	r.Use(middleware.IpLimit())
}

func loadRouter(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		//==============对象=================

		//创建对象
		v1.POST("/classes/:className", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("create", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := q.Create()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusCreated, *res)
				}
			}
		})

		//获取对象
		v1.GET("/classes/:className/:objectId", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("get", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := q.Get()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusCreated, *res)
				}
			}
		})

		//更新对象
		v1.PUT("/classes/:className/:objectId", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("update", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := q.Update()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusCreated, *res)
				}
			}
		})

		//查询对象
		v1.GET("/classes/:className", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("find", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := q.Find()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusCreated, gin.H{"results": *res})
				}
			}
		})

		//删除对象
		v1.DELETE("/classes/:className/:objectId", func(c *gin.Context) {
			cont := NewContext(c)
			q, err := GenQuery("delete", c, cont)
			if !cont.ReturnError(err) {
				err = Authenticate(cont)
			}
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				err = q.Delete()
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, gin.H{"success": true})
				}
			}
		})

		v1.POST("/newClass", func(c *gin.Context) {
			//cont:=NewContext(c)
			CreateTable(c.PostForm("className"), map[string]interface{}{})
		})

		v1.POST("/newfield", func(c *gin.Context) {
			//cont:=NewContext(c)
			CreateField(c.PostForm("className"), c.PostForm("fieldName"), c.PostForm("type"), c.DefaultPostForm("default", ``), false, false, false, false, c.DefaultPostForm("notes", ``), c.PostForm("relationTo"), false)
		})

		//==============用户=================

		//用户注册 用户连接=
		v1.POST("/users", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.UserRegister(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//用户登录=
		v1.GET("/login", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.UserLogin(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//重置用户 sessionToken
		v1.PUT("/users/:objectId/refreshSessionToken	", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.ResetSessionToken(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//更新密码，要求输入旧密码
		v1.PUT("/users/:objectId/updatePassword", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.ResetPassword(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//请求密码重设============
		v1.POST("/requestPasswordReset", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.RequestPasswordReset(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//邮箱-密码重设
		v1.GET("/passwordReset/:email/:yzm", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.EmailPasswordReset(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//请求验证用户邮箱
		v1.POST("/requestEmailVerify", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.RequestEmailVerify(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//验证用户邮箱
		v1.GET("/emailVerify/:email/:yzm", func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibuser.EmailVerify(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//========================文件==========================

		//获得上传凭证
		v1.GET(`/files/token`, func(c *gin.Context) {
			cont := NewContext(c)
			q := ibfile.GetUploadQuery(cont)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				c.JSON(http.StatusOK, ibfile.UploadFile(cont))
			}
		})

		//七牛回调接口
		v1.POST(`/files/callback`, func(c *gin.Context) {
			cont := NewContext(c)
			res, err := ibfile.UploadCallback(cont)
			if !cont.ReturnError(err) {
				c.JSON(http.StatusOK, res)
			}

		})

		//删除文件
		v1.DELETE(`/files/:objectId`, func(c *gin.Context) {
			cont := NewContext(c)
			q := ibfile.GetDeleteQuery(cont)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				err = AclFilter(cont, q)
			}
			if !cont.ReturnError(err) {
				res, err := ibfile.DeleteFile(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}

		})

		//========================master==========================

		//测试返回

		v1.GET(`/master/test`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				if cont.IsMaster {
					c.JSON(http.StatusOK, gin.H{"success": true, "msg": "master connect"})
				} else {
					c.JSON(http.StatusOK, gin.H{"success": false})
				}
			}
		})

		//创建Class
		v1.POST(`/master/class`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.CreateClass(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//查询所有Class
		v1.GET(`/master/allClass`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetAllClass(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//设置Class权限
		v1.PUT(`/master/class/:className`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.SetClassAcl(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//删除Class
		v1.DELETE(`/master/class/:className`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.DeleteClass(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//删除Class所有数据
		v1.DELETE(`/master/allData/:className`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.DeleteClassAllData(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//获取Class字段信息
		v1.GET(`/master/field/:className/:fieldName`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//获取Class所有字段信息
		v1.GET(`/master/allFieldInClass/:className`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetFieldInClass(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//获取所有字段信息
		v1.GET(`/master/allField`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetAllField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//创建字段
		v1.POST(`/master/field/:className/`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.CreateField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//删除字段
		v1.DELETE(`/master/field/:className/:fieldName`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.DeleteField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//修改字段
		v1.PUT(`/master/field/:className/:fieldName`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.UpdateField(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//获取应用属性
		v1.GET(`/master/appInfo`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.GetAppInfo(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//修改应用属性
		v1.PUT(`/master/appInfo/:key`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmaster.SetAppInfo(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

		//=================系统监控===================

		//获取系统资源使用情况
		v1.GET(`/monitor/systemInfo`, func(c *gin.Context) {
			cont := NewContext(c)
			err := Authenticate(cont)
			if !cont.ReturnError(err) {
				res, err := ibmonitor.GetSystemInfo(cont)
				if !cont.ReturnError(err) {
					c.JSON(http.StatusOK, res)
				}
			}
		})

	}
}
