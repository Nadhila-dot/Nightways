# AI education enviroment.

Artifically inteligent Education Enviroment.


Strucutre

A sheet is the basic structure of an assigment.
It's as said a sheet, a notebook can have multiple sheets.
A sheet has
- Title
- Descrption
- Name

Notebooks are subject based sheets
Each notebook has 1 or many sheets related to a subject.
A notebook has
- Name
- tags array[] (the course keywords and etc)
- Title
- Description

A note book has a stack.json to show where it came from, the difculity of each sheet and etc. Sort of like the metadata

each sheet has a hidden base64 encoded string on top prob in whitetext or black text called id
Id holds the 
source ai : <gemini, grok and etc>
Date: <date it was created>
Misc: <misc data>


Queue System.

FiFo

First in First out. 

Possion checkers, it will check if something failed, then notify the frontend and then start to remove itself from the database. 



Remeber
in the backend websockets
type push = simple send message to user
type update = proper update message with data



TODO
- Fix notifications system (done)
- Fix qeue (done)
- Export metadata to the bucket (not done yet)
- Add a workflow where in generation you can step in and change things (not done yet)
- Make it so that you can the latex and regenerate stuff (done)
- Fix the json updates sent via the update. (kinda done)
- Fix desktop and Mobile bundles



