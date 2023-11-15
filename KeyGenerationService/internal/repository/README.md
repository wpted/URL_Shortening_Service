## Two choices for storing keys and used_keys.
### Create another table for used keys
    Pros:
	    1. Separation of Concerns - Maintaining data and querying is easier
	    2. Simplified Queries - No filtering
	    3. Easier Archiving or Deletion - Modification made to table is easier.
    Cons:
	    1. Complexity - Maintaining multiple table increased complexity in database schema. 
### Add another column to mark keys as used
    Pros:
        1. Simplicity - Adding column to mark keys as used in the same table can be a simpler, straightforward design.
        2. Queries - Queries can involve both active and used keys.
    Cons:
        1. Increased Table Size- Could impact performance
   
