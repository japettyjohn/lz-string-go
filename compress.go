package lzstring

func compress(uncompressed string) ([]byte,error) {
    var (
          dictionary=map[string]int64{}
          wc, w, result string
          enlargeIn= 2 // Compensate for the first entry which should not count
          dictSize= 3
          numBits= 2
          data= struct {
          	d string
          	val int64
          	position int64
          }{}
        	
        )
    for _,c:= range uncompressed {
      if (!Object.prototype.hasOwnProperty.call(context.dictionary,context.c)) {
        context.dictionary[context.c] = context.dictSize++;
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
}
