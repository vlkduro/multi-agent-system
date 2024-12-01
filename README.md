# AI30 A24 - Projet - Valakou Martins Chartier Bidaux

| Information | Valeur                                                                         |
| ----------- | ------------------------------------------------------------------------------ |
| Auteurs     | Valakou Quentin <br> Martins Clément <br> Chartier Julie <br> Bidaux Alexandre |
| Date        | 27/11/24                                                                       |
| Version     | 0.0.1                                                                          |

## Sommaire

- [Architecture actuelle](#architecture-actuelle)
- [Conception](#conception)
- [Sources](#sources)

## Architecture actuelle

`Agent` et `Object` sont des interfaces de même stabilité que l'environement. C'est à dire que lorsque une des interfaces change, `Environment` change également. Ils sont donc dans le même package pour l'instant

> Remarque : mettre `Agent` (i.e `Object`) dans le package `agent` (i.e `object`) crée une boucle de dépendance :
> `Agent` appelle le package `environment` et `Environment` appelle le package `agent`, par défaut le compilateur interdit les dépendances cycliques

## Conception

### Problématique

Quelles proportions d’abeilles faut-il à une colonie pour prospérer face au danger des frelons ? (_Apis Mellifera_ (abeille européenne à miel) vs _Vespa crabo_ (frelons européeens) ou _Vespa velutina_ (frelons asiatiques) apparemment les abeilles françaises ne savent pas encore s'En défendre, à vérifier) (tentative)

### Courte description

### Règles

- Une abeille vit dans une ruche
- Une abeille trouve des fleurs à butiner
- Une abeille avec sa capacité de pollen au maximum rentre déposer le pollen à la ruche
- Une abeille qui rencontre une fleur pendant son retour retiens sa position
- Une abeille connaissant la position d’au moins une fleur partage ses connaissances aux autres abeilles
- Les abeilles fuient les frelons
- La ruche contient une reine qui produit des abeilles en consommant du miel
- La ruche contient des abeilles productrices qui transforment le pollen en miel
- La ruche meurt si elle n’a pas suffisamment de miel
- La ruche peut être déplacée contre une grosse quantité de miel (proportionnel au nombre d’habitants)
- Une ruche peut tuer un frelon si il y a suffisamment d’abeilles à l'intérieur (trouver le ratio en ligne)
- Un frelon chasse les abeilles
- Un frelon cherche les ruches d’abeilles
- Suffisamment de frelons peuvent tuer la reine
- Un frelon peut environner des fleurs, suivre des abeilles ou environner la ruche

On pourrait proposer à l'utilisateur de gérer la proportion d'abeilles productrices et exploratrices en temps réel (chaque abeille ayant un coût en miel pour la ruche).

## Sources

### Danse des abeilles

- [« The Waggle Dance of the Honeybee »](https://www.youtube.com/watch?v=bFDGPgXtK-U), Georgia Tech College of Computing, 03/02/2011 (visionnée le 28/11/2024).
- [« Bee Dance Language - the linguistics behind animal language »](https://www.youtube.com/watch?v=pb1lRI-YePU), NativLang, 20/12/2014 (visionnée le 28/11/2024)
- [« Les moyens de communication impressionnants des abeilles »](https://www.nationalgeographic.fr/video/animaux/les-moyens-de-communication-impressionnants-des-abeilles), National Geographic, 14/09/2023 (visionnée le 28/11/2024)

### Généralités et apiculture

- Podcast Le bruit de la ruche (Saison 1, épisodes 1 à 5 écoutés le 28/11/2024)
- Documentaire « Des abeilles et des hommes - Et si les abeilles disparaissaient », Markus Imhoof, 2012 (visionné le 28/11/2024)

### Prédateurs des abeilles

- https://www.apiculteurs-occitanie.fr/informations/predateurs/ (consulté le 28/11/2024) : prédateurs des abeilles

### Défense contre les prédateurs

- [« Le redoutable système de défense de l'abeille japonaise face à un frelon »](https://www.youtube.com/watch?v=bDHWv8_-f4U), National Geographic Wild France, 21/11/2020(visionnée le 28/11/2024)
  - les abeilles japonaises sont capables de tuer un frelon, dès sa première victime les abeilles se mettent sur le frelon et agitent leur dard toutes ensemble pour créer de la chaleur. Les abeilles japonaises supportent jusqu’à 48 degré celsius alors que le frelon ne survit pas à plus de 46 degré (quel type de frelon ?)
- [« Ce frelon n'aurait pas dû se frotter à ces abeilles »](https://www.youtube.com/watch?v=N_AO8XvAZgQ), National Geographic Wild France, 28/04/2021 (visionnée le 28/11/2024)
  - les gardiennes se jettent sur le frelon. Les abeilles sont programmées pour défendre la colonie même si elles doivent y laisser la vie. Piqûres d’abeilles ne tuent pas les frelons mais les affaiblissent. Les frelons fuient mais le dernier est cuit de la même manière que dans la précédente vidéo. Pas d’info sur les espèces exactes.
- [« Comment les abeilles parviennent-elles à se défendre ? »](https://www.nationalgeographic.fr/video/animaux/comment-les-abeilles-parviennent-elles-a-se-defendre), National Geographic, 15/04/2023 (visionnée le 28/11/2024)
  - Api Mellifica (abeilles à miel, sociales, les plus connues), pollinisatrices
    - Seules les femelles ont un dard en bas de leur abdomen.
    - Cousines des fourmis et des guèpes
    - Essaim : colonie d’Abeilles qui forme un super-organisme, indivisible. Se crée quand la reine et quelques abeilles quittent la ruche pour chercher un endroit pour former une nouvelle ruche
    - Phéromone d’alerte propagé dans la colonie ? et attaquent pour protéger la reine et la colonie
