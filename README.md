# AI30 A24 - Projet - Valakou Martins Chartier Bidaux

|Information|Valeur|
|-|-|
|Auteurs|Valakou Quentin|
||Martins Clément|
||Chartier Julie|
||Bidaux Alexandre|
|Date|27/11/24|
|Version|0.0.1|

## Architecture actuelle 

`Agent` et `Object` sont des interfaces de même stabilité que l'environement. C'est à dire que lorsque une des interfaces change, `Environment` change également. Ils sont donc dans le même package pour l'instant 

> Remarque : mettre `Agent` (i.e `Object`) dans le package `agent` (i.e `object`) crée une boucle de dépendance :
> `Agent` appelle le package `environment` et `Environment` appelle le package `agent`, par défaut le compilateur interdit les dépendances cycliques
