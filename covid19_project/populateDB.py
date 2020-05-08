import databaseconfig as cfg
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
        db_Info = connection.get_server_info()
        print("Connected to MySQL Server version ", db_Info)
        cursor = connection.cursor()
        cursor.execute("select database();")
        record = cursor.fetchone()
        print("You're connected to database: ", record)

except Error as e:
    print("Error while connecting to MySQL", e)


##----------------DATABASE CONNECTED-----------------------------------------
    
finally:
    if (connection.is_connected()):




        ##---------------------CREATING TABLE IF NOT PRESENT---------------
        
        
        cursor.execute("SHOW TABLES LIKE '"+cfg.mysql["table_name"]+"';")
        result = cursor.fetchone()
        if result:
            # there are is a table named "tableName"
            print("given table already exist droping and creating again")
            cursor.execute("drop table "+cfg.mysql["table_name"]+";")
            create_table="""CREATE TABLE CountryWiseData
                            (CountryName varchar(255),
                             ProvinceName varchar(255),
                             ActiveCases varchar(255),
                             TotalRecovered int,
                             TotalDeaths int,
                             Date date)"""
            cursor.execute(create_table)
        else:
            # there are no tables named "tableName"
            print("given table does not exist creating now")
            create_table="""CREATE TABLE CountryWiseData
                            (CountryName varchar(255),
                             ProvinceName varchar(255),
                             ActiveCases varchar(255),
                             TotalRecovered int,
                             TotalDeaths int,
                             Date date)"""
            cursor.execute(create_table)





        ##--------------------INSERTING DATA-------------------------------
            
            
        print("Inserting Data")
        insert_query = """INSERT INTO CountryWiseData 
                          (CountryName,ProvinceName,ActiveCases,
                           TotalRecovered,TotalDeaths,Date)
                           VALUES (%s, %s, %s, %s, %s, %s)"""
        
        count=0  #-----FOR COUNTING NO. OF NEWLY INSERTED DATA--------------
    
        for x in data:
            country=x['Country']
            province=x['Province']
            city=x['City']
            date=(x['Date'])
            date=date.split('T')
            death=str(int(x['Deaths']))
            confirmed=str(int(x['Confirmed']))
            recovered=str(int(x['Recovered']))
            val = (country,province,confirmed,recovered,death,date[0])
            cursor.execute(insert_query, val)
            count=count+1

        print("No. Of Rows Affeced:%s"%(count))
        connection.commit()

        cursor.close()
        connection.close()
        print("MySQL connection is closed")
