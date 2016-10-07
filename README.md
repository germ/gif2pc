Gif -> Point Cloud
--

A simple utility to visualize videos as solid 3d objects.
Main.go takes one argument, a gif, and outputs a series of
scripts and point clouds. Each point cloud represents a color
in the original gif (for preformance reasons, use simple color
gifs for testing). These scripts are run against the clouds using
Meshlab to try aproximating their shape and OBJ files are placed
in a output directory. These OBJ files can then be loaded into
blender for animating/exploration.

Workflow
1| go run ../main.go globe.gif
2| ../aaRUNME.sh
3| Import out/*.obj to Blender

This is a toy, use it as you will.
Demo Video: https://www.instagram.com/p/BLQUKEyAD3T/
