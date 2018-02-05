# AntColony
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
- visitQuantity: The number of towns that should be visited over the duration of the trip

### Towns
- includesHome: True or false indicating whether the first town in the provided array is a "home" location where each ant should start and return to every time. If set to false, there will be no set starting location and the starting town will be set randomly.
- towns: the array of towns with all of the relevant properties associated as follows
  - id: Unique id associated with the town
  - distances[]: Distances to each other location in the object. The order is assumed to be identical as the `towns` array, and the algorithm will return an error if the length of `distances` is not equal to `towns`
  - trails[]: (Optional) The pheremone trail associated with the path to the other towns in the case that the results are being sent back to the algorithm for more iterations.
  - rating: The desirability factor of the town. when `config.maximizeRating == true`, higher ratings will increase the probability that an ant will visit this town
  - isRequired: Marks the town as mandatory - If set to true while `config.visitQuantity` is less than the size of `towns.towns[]`, any ant that fails to visit this town will be considered invalid and another ant will take its place.

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