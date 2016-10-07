
#!/bin/bash
mkdir out
len=`ls -l | grep txt | wc -l`
for i in `seq 0 $((len-1))`; do meshlabserver -i $i.txt -s shape.mlx -om fc -o $i.obj; meshlabserver -i $i.obj -s $i.mlx -om fc -o out/$i.obj;  done
