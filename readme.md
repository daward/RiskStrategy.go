This is an analysis of the Risk board game played on an approximate (if not dated) [world map](https://i.pinimg.com/564x/ae/44/2d/ae442d7e848ad1bce549d2f21ee73f9c.jpg). The intention is to provide a player insight of relative country values and a guide to understand how they should conduct their world war.

# Method
Like most map-based games, Risk can be thought of as an undirected graph for the purposes of analysis. The algorithm uses a breadth first search to flood the map with various conquest permutations evaluating each position hueristically. The heuristic is an accumulating score that is defined as `reinforcements \ border defense territories`. 

For example, if a player were to start the game owning the whole Australian continent, they will receive 3 reinforcements but only have to defend 1 territory. They will have a total cumulative score at that time of 3. If they choose to conquer Siam next, they will still have a score of 3, which will be added to their previous score, resulting in a total of 6. Each option is evaluated (within a pruning tolerance) and cumulative scores are compared.

# Results
The complete results can be viewed at [results.json](results.json). The results are generated considering exactly one territory (there are 42 in all) as a starting position and running the algorithm to get a final score for that single country starting position.

# Conclusion
My hypothesis was that the Australian start would easily outclass others due to its hot start and dominating 3-score position that is unmatched on the board. The counter argument appears to be true however. It is hard to break out of Australia, and the player is best served to conquer Asia. During this time, other starting positions have caught up to the Australian start, particularly those in North America. Indeed, with the exception of Northern Europe, North American starts (especially those in the east) dominate the top positions. The middle east also offers a decent option with its ability to enter Africa or protect Europe.