import databaseconfig as cfg
from datetime import datetime as dt
import json
import requests
from dictor import dictor


##--------------FETCHING DATA FROM API----------------------------------

URL = "https://api.covid19api.com/all"

r = requests.get(url = URL) 
 
data = r.json() 

import mysql.connector
from mysql.connector import Error

##--------------CREATING DATABASE IF NOT PRESENT------------------------

try:
    connection = mysql.connector.connect(host=cfg.mysql["host"],
                                         user=cfg.mysql["user"],
                                         password=cfg.mysql["password"])
    if connection.is_connected():
        db_Info = connection.get_server_info()
        print("Connected to MySQL Server version ", db_Info)
        cursor = connection.cursor()
        cursor.execute("SHOW DATABASES LIKE '"+cfg.mysql["db"]+"';")
        result = cursor.fetchone()
        if result:
            # there is a database named "databaseName"
            print("database already present")
            
        else:
            # there are no databases named "databaseName"
            print("database not present!!  Creating database now")
            cursor.execute("CREATE DATABASE "+cfg.mysql["db"]+";")

except Error as e:
    print("Error while connecting to MySQL", e)


##---------------CONNECTING TO THE DATABASE-----------------------------
    
try:
    connection = mysql.connector.connect(host=cfg.mysql["host"],
                                         database=cfg.mysql["db"],
                                         user=cfg.mysql["user"],
                                         password=cfg.mysql["password"])
    if connection.is_connected():
        cursor = connection.cursor()
        cursor.execute("select database();")
        record = cursor.fetchone()
        print("You're connected to database: ", record[0])

except Error as e:
    print("Error while connecting to MySQL", e)


##----------------DATABASE CONNECTED-----------------------------------------
    
finally:
    if (connection.is_connected()):




        ##---------------------CREATING TABLE IF NOT PRESENT---------------

              
        cursor.execute("SHOW TABLES LIKE '"+cfg.mysql["table_name"]+"';")
        result = cursor.fetchone()
        if result:
            ##----------------TABLE IS ALREADY PRESENT---------------------
            
            print("given table already exist ");
            print("Connected to table:",cfg.mysql["table_name"])

            ##----------------FINDING LATEST DATA--------------------------
            
            max_query="""SELECT MAX(Date) FROM """+cfg.mysql["table_name"]+""";"""
            cursor.execute(max_query)
            last_date = cursor.fetchone()
            if last_date[0] is None:
                my_string= "2016-06-11"
                last_date = (dt.strptime(my_string, "%Y-%m-%d").date(),)
            else:
                dateStr = last_date[0].strftime("%Y-%m-%d ")           
                print('Latest Date : ' ,dateStr)
            
        else:
            ##-----------TABLE IS NOT PRESENT!! CREATING TABLE---------------

            print("given table does not exist");
            create_table="""CREATE TABLE """+cfg.mysql["table_name"]+"""
                            (CountryName varchar(255),
                             ProvinceName varchar(255),
                             ActiveCases varchar(255),
                             TotalRecovered int,
                             TotalDeaths int,
                             Date date)"""
            
            cursor.execute(create_table)
            print("Connected to table:",cfg.mysql["table_name"])

            ##----SETTING AN OLD DATE IN LAST_DATE FOR FIRS TIME ENTRY------
            
            my_string= "2016-06-11"
            last_date = (dt.strptime(my_string, "%Y-%m-%d").date(),)




        ##--------------------INSERTING DATA-------------------------------
            
            
        print("Inserting Data")
        insert_query = """INSERT INTO """+cfg.mysql["table_name"]+""" 
                          (CountryName,ProvinceName,ActiveCases,
                           TotalRecovered,TotalDeaths,Date)
                           VALUES (%s, %s, %s, %s, %s, %s)"""
        
        count=0  #-----FOR COUNTING NO. OF NEWLY INSERTED DATA--------------
    
        for x in data:
            
            date=(x['Date'])
            date=date.split('T')
            if date[0]>last_date[0].strftime("%Y-%m-%d"):
                country=x['Country']
                province=x['Province']
                city=x['City']
                death=str(int(x['Deaths']))
                confirmed=str(int(x['Confirmed']))
                recovered=str(int(x['Recovered']))
                val = (country,province,confirmed,recovered,death,date[0])
                cursor.execute(insert_query, val)
                count=count+1

                
        print("No. Of Rows Affeced:",count)

        connection.commit()

        cursor.close()
        connection.close()
        print("MySQL connection is closed")



    
    
