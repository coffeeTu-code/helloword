import os
import sys
import io
import json
import argparse

import base64
#import imageio
import cv2
import numpy as np
import scipy.io as sio
import multiprocessing as mp

from PIL import Image
from scipy import misc

def process(sample, source_part, target_part, symbol_part, lowpass=13):

        src_png = "{}/{}.png".format("debug", sample['id'])

        w = int(sample['meta']['bandwidth'])
        h = int(sample['meta']['features'])
        counter = np.zeros((4,), dtype=np.int32)

        src = np.zeros( (h, w), dtype=np.uint8)
        y = 0

        counter[1] += h*w
        
        for sdr in sample['sdr']:

            x_ = sdr[0]
            z_ = sdr[1]

            n = len(sdr[0]) if len(sdr[0]) < lowpass else lowpass

            for i in range(n):
                x = int(x_[i])
                z = int(z_[i])
                src[y, x] = z
                counter[2] += 1
                counter[3] += z

            y = y+1

        #imageio.imwrite(src_png, src)
        is_success, im_buf_arr = cv2.imencode(".png", src)

        assert (is_success), "cannot encode png file"
        byte_im = im_buf_arr.tobytes()
        source_part.append(base64.b64encode(byte_im))

        tgt = np.zeros( (1,2,1), dtype=np.uint8)

        tgt[0,0,0] = 255*(1 - sample['target'])
        tgt[0,1,0] = 255*sample['target']

        #imageio.imwrite(tgt_png, tgt)
        is_success, im_buf_arr = cv2.imencode(".png", tgt)
        assert (is_success), "cannot encode png file"
        byte_im = im_buf_arr.tobytes()
        target_part.append(base64.b64encode(byte_im))

        symbol_part.append(sample['id'])

        if(sample['target'] == 1):
            counter[0] += 1
        
        return counter

parser = argparse.ArgumentParser(description='')
parser.add_argument('--source', dest='source', default="sample.jss", help='path for data source')
parser.add_argument('--output', dest='output', default="asset", help='path for data source')
parser.add_argument('--prefix', dest='prefix', default="train", help='prefix for partitions')
parser.add_argument('--cutoff', dest='cutoff', type=int, default=-1, help='# of observations')
parser.add_argument('--lowpass', dest='lowpass', type=int, default=13, help='# of observations')
args = parser.parse_args()

def make_dir(path):
  if not os.path.exists(path):
    os.makedirs(path)

def parts_process(packaged):

    parts, offset = packaged

    source_part = []
    target_part = []
    symbol_part = []

    sums = np.zeros((4,), dtype=np.int32)

    for line in parts:
        encoded = json.loads(line)
        counter = process(encoded, source_part, target_part, symbol_part, lowpass=args.lowpass)
        sums = sums + counter
    
    #print()
    #print(source_part[0])
    #print()
    #print(target_part[0])

    start = offset - len(parts)
    partition_file = "{}/{}-{}+{}.part".format(args.output, args.prefix, start, len(parts))
    sio.savemat(partition_file, { "symbols":symbol_part, "targets":target_part, "sources":source_part }, do_compression=True )
    return offset, sums


def main():

    filepath = args.source
    cutoff = args.cutoff

    if not os.path.isfile(filepath):
       print("File path {} does not exist. Exiting...".format(filepath))
       sys.exit()

    asset_dir = args.output
    make_dir(asset_dir)

    worker = mp.cpu_count()
    package = 1024
  
    with open(filepath) as fp:
        cnt = 0
        sums = np.zeros((4,), dtype=np.int64)

        partitions = []
        packaged = []

        lcnt = 0

        for line in fp:
            lcnt += 1
            packaged.append(line)
            if(len(packaged) >= package):
                partitions.append( (packaged,lcnt) )
                packaged = []

        if(len(packaged)):
            partitions.append( (packaged,lcnt) )

        print(len(packaged))

        print("PARTITION {} -> {}".format(len(partitions), lcnt))

        #parts_process(partitions[-1])
        #exit()

        with mp.Pool(worker) as pool:
            for x in pool.imap(parts_process, (x for x in partitions), 1):
                cnt, counter = x
                sums = sums + counter
                print(cnt, " processed, positive ratio={}, sparsity={}, mean={}".format(sums[0]/cnt, sums[2]/sums[1], sums[3]/sums[2]))

        #print("________DONE__________")
        #print(cnt, " processed, positive ratio={}, cardinal={}, lowpass ratio={}".format(sums[0]/cnt, sums[2]/sums[1], sums[2]/sums[3]))

if __name__ == '__main__':
   main()