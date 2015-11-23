package encoding

import (
	"fmt"
	"math"
)

/*
 var context = {
      dictionary: {},
      dictionaryToCreate: {},
      c:"",
      wc:"",
      w:"",
      enlargeIn: 2, // Compensate for the first entry which should not count
      dictSize: 3,
      numBits: 2,
      result: "",
      data: {string:"", val:0, position:0}
    }, i;

    for (i = 0; i < uncompressed.length; i += 1) {
      context.c = uncompressed.charAt(i);
      if (!Object.prototype.hasOwnProperty.call(context.dictionary,context.c)) {
        context.dictionary[context.c] = context.dictSize++;
        context.dictionaryToCreate[context.c] = true;
      }

      context.wc = context.w + context.c;
      if (Object.prototype.hasOwnProperty.call(context.dictionary,context.wc)) {
        context.w = context.wc;
      } else {
        this.produceW(context);
        // Add wc to the dictionary.
        context.dictionary[context.wc] = context.dictSize++;
        context.w = String(context.c);
      }
    }

    // Output the code for w.
    if (context.w !== "") {
      this.produceW(context);
    }

    // Mark the end of the stream
    this.writeBits(context.numBits, 2, context.data);

    // Flush the last char
    while (context.data.val>0) this.writeBit(0,context.data)
    return context.data.string;
*/

type compressor struct {
	dictionary         map[rune]int32
	dictionaryToCreate map[rune]interface{}
	wc                 rune
	w                  rune
	result             string
	enlargeIn          int
	dictSize           int32
	numBits            int
	data               string
	val                int32
	position           int
}

func newCompressor() *compressor {
	return &compressor{
		dictionary:         map[rune]int32{},
		dictionaryToCreate: map[rune]interface{}{},
		enlargeIn:          2,
		dictSize:           3,
		numBits:            2,
	}
}

func (c *compressor) compress(uncompressed string) (string, error) {
	for _, char := range uncompressed {
		if _, ok := c.dictionary[char]; !ok {
			c.dictionary[char] = c.dictSize
			c.dictSize++
			c.dictionaryToCreate[char] = nil
		}

		c.wc = c.w + char
		if _, ok := c.dictionary[c.wc]; ok {
			c.w = c.wc
		} else {
			// productW
			c.dictionary[c.wc] = c.dictSize
			c.dictSize++
			c.w = char
		}
	}

	if c.w != 0 {
		c.produceW()
	}

	c.writeBits(c.numBits, 2)

	// Flush the last char
	for c.val > 0 {
		c.writeBit(0)
	}
	return c.data, nil
}

/*
  writeBit : function(value, data) {
    data.val = (data.val << 1) | value;
    if (data.position == 15) {
      data.position = 0;
      data.string += String.fromCharCode(data.val);
      data.val = 0;
    } else {
      data.position++;
    }
  },

  writeBits : function(numBits, value, data) {
    if (typeof(value)=="string")
      value = value.charCodeAt(0);
    for (var i=0 ; i<numBits ; i++) {
      this.writeBit(value&1, data);
      value = value >> 1;
    }
  }
*/
func (c *compressor) writeBit(value rune) {
	c.val = (c.val << 1) | value
	if c.position == 15 {
		c.position = 0
		c.data = fmt.Sprintf("%s%c", c.data, c.val)
		c.val = 0
	} else {
		c.position++
	}
}

func (c *compressor) writeBits(numBits int, value interface{}) {
	var finalValue rune
	switch t := value.(type) {
	case string:
		finalValue = []rune(t)[0]
	case rune:
		finalValue = t
	}
	for i := 0; i < numBits; i++ {
		c.writeBit(finalValue & 1)
		finalValue = finalValue >> 1
	}
}

/*

  produceW : function (context) {
    if (Object.prototype.hasOwnProperty.call(context.dictionaryToCreate,context.w)) {
      if (context.w.charCodeAt(0)<256) {
        this.writeBits(context.numBits, 0, context.data);
        this.writeBits(8, context.w, context.data);
      } else {
        this.writeBits(context.numBits, 1, context.data);
        this.writeBits(16, context.w, context.data);
      }
      this.decrementEnlargeIn(context);
      delete context.dictionaryToCreate[context.w];
    } else {
      this.writeBits(context.numBits, context.dictionary[context.w], context.data);
    }
    this.decrementEnlargeIn(context);
  },

  decrementEnlargeIn : function(context) {
    context.enlargeIn--;
    if (context.enlargeIn == 0) {
      context.enlargeIn = Math.pow(2, context.numBits);
      context.numBits++;
    }
  },

*/

func (c *compressor) produceW() {

	if _, ok := c.dictionaryToCreate[c.w]; ok {
		if c.w < 256 {
			c.writeBits(c.numBits, 0)
			c.writeBits(8, c.w)
		} else {
			c.writeBits(c.numBits, 1)
			c.writeBits(16, c.w)
		}
		c.decrementEnlargeIn()
		delete(c.dictionaryToCreate, c.w)
	} else {
		c.writeBits(c.numBits, c.dictionary[c.w])
	}
	c.decrementEnlargeIn()
}

func (c *compressor) decrementEnlargeIn() {
	c.enlargeIn--
	if c.enlargeIn == 0 {
		c.enlargeIn = int(math.Pow(2, float64(c.numBits)))
		c.numBits++
	}
}
