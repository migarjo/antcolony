# AntColony

[![CircleCI](https://circleci.com/gh/migarjo/antcolony.svg?style=svg)](https://circleci.com/gh/migarjo/antcolony)

This is an [ant colony optimization](http://iridia.ulb.ac.be/~mbiro/Paperi/IridiaTr2006-023r001.pdf) implementation with the purpose of solving the [Travelling Salesman Problem](https://en.wikipedia.org/wiki/Travelling_salesman_problem).

## Using AntColony

It can be accessed at https://migarjo-ant-colony.herokuapp.com/api/solvetsp

This service can be used at the above endpoint with a POST call, comprising of `config` and `towns` objects with the following properties defined:

### Config
- numberOfIterations: The number of times a group of ants will traverse the towns provided
- ratingPreference: The relative weight given to the rating, or desirability associated with a particular town (see `towns.towns[].rating`)
- maximizeRating: True or false indicating whether higher rating is desirable - if `maximizeRating == true`, the higher rating will be preferred over the lower.
- distancePreference: The relative weight given to minimizing the distance traveled between two towns (see `towns.towns[].distances[]`)
- trailPreference: The relative weight given to pheremone trails deposited by ants from previous iterations
  - _Note:_ if this is set to 0, the algorithm with be completely random, with no average improvement of ants over multiple iterations
- pherememoneStrength: The baseline strength of the pheremone trail deposited by each ant
- evaporationRate: The ratio at which the pheremones of ants from previous iterations decreases over time
- visitQuantity: The maximum number of towns that should be visited over the duration of the trip
- verbose: If set to `true`, additional history will be processed in the algorithm and returned to the requester in the REST call
- tripBounds: Sets the start and end bounds of the trip. If `tripBounds.start > 0`, the trip will start counting from that point. If `tripBounds.end > 0`, the trip will end if no remaining unvisited towns can be visited without exceeding the given value.
- minimumTripUsage: Sets a default ratio of the total time in `tripBounds` that each ant must consume before the trip ends. If an ant takes a path that eliminates any options to reach the minimum trip usage, it will reattempt to take another path until it finds one. If it reaches 100 attempts without meeting the minimum trip length, the trip will iteratively be reduced by a factor of 0.9, until the ant can complete the trip.

### Towns
- includesHome: True or false indicating whether the first town in the provided array is a "home" location where each ant should start and return to every time. If set to false, there will be no set starting location and the starting town will be set randomly.
- towns: The array of towns with all of the relevant properties associated as follows
  - id: Unique id associated with the town
  - distances[]: Distances to each other location in the object. The order is assumed to be identical as the `towns` array, and the algorithm will return an error if the length of `distances` is not equal to `towns`
  - trails[]: (Optional) The pheremone trail associated with the path to the other towns in the case that the results are being sent back to the algorithm for more iterations.
  - rating: The desirability factor of the town. when `config.maximizeRating == true`, higher ratings will increase the probability that an ant will visit this town
  - isRequired: Marks the town as mandatory - If set to true while `config.visitQuantity` is less than the size of `towns.towns[]`, any ant that fails to visit this town will be considered invalid and another ant will take its place.
  - trailHistory: Returns the trail results for each town at each iteration if `config.verbose == true`.
  

  ### Sample request body
  The following is a sample request body consisting of five points arranged in a regular pentagon.

  ```json
{ 
	"config":  {
        "numberOfIterations": 10,
        "trailPreference": 1,
        "ratingPreference": 1,
        "distancePreference": 1,
        "pherememoneStrength": 1,
        "evaporationRate": 0.8,
        "maximizeRating": true,
        "visitQuantity": 4
	},
	"towns":{
        "includesHome": true,
		"towns": [
			{
				"id": 0,
				"distances": [
					0,
					1,
					1.618,
                    1.618,
                    1
				],
				"trails": [
					0.1,
					0.2,
					0.3,
                    0.4,
                    0.5
				],
				"rating": 1
			},
			{
				"id": 1,
				"distances": [
					1,
					0,
					1,
                    1.618,
                    1.618
				],
				"trails": [
					0.1,
					0.2,
					0.3,
                    0.4,
                    0.5
				],
				"rating": 1
			},
			{
				"id": 2,
				"distances": [
					1.618,
					1,
					0,
                    1,
                    1.618
				],
				"trails": [
					0.1,
					0.2,
					0.3,
                    0.4,
                    0.5
				],
				"rating": 1
			},
			{
				"id": 3,
				"distances": [
					1.618,
					1.618,
					1,
                    0,
                    1
				],
				"trails": [
					0.1,
					0.2,
					0.3,
                    0.4,
                    0.5
				],
				"rating": 1
			},
			{
				"id": 4,
				"distances": [
					1,
					1.618,
					1.618,
                    1,
                    0
				],
				"trails": [
					0.1,
					0.2,
					0.3,
                    0.4,
                    0.5
				],
				"rating": 1
			}
		]
	}
}
```