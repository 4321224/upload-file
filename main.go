package main

import (
    "context"
    "io"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Replace with your MongoDB server URI
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	
	// Access the "fileuploads" database and the "files" collection
	db := client.Database("fileuploads")
	collection := db.Collection("files")

	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

        src, _ := file.Open()
        defer src.Close()
	
	    dst, err := os.Create("uploads/" + file.Filename)
	    if err != nil {
		    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		    return
	    }
        defer dst.Close()
	
	    if _, err = io.Copy(dst, src); 	err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": 	err.Error()})
            return
        }
	
	    // Insert information about the uploaded file into MongoDB
	    fileInfo := bson.M{
		    "filename": file.Filename,
		    "filepath": 	"uploads/" + file.Filename,
		    // Add any other relevant metadata about the file
	    }

	    _, insertErr := collection.InsertOne(context.Background(), fileInfo)
	    if insertErr != nil {
		    c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		    return 
	    }

        c.JSON(http.StatusOK, gin.H{
           	  	 	 	 	 	          	       	       	    	 	 	  	    	  	  	             	              	           	           	            	             	              	         	      	     	          	        	    	 	                	      	             	              	         	      	     	          	        	    	 	                	      	             	              	         	      	     	          	        	    	 	                 		 		 	    		 			 		 				 				 					 					  	        	  			 	  		  		  			 			   			     		      			      			  			    			  				  			  				   		     			     			       		       		       		        		                  		           		           		            		                  		           		            		                  		           		            		                  		           		            		                  		           		            		                  		           		                   		            		               		            		                   		            		               		            		                   		            		               		         		        		    	 		                 	                 	        	        	  	            	  			 	  		  		  			 			   			     		      			   
