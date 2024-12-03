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

### Architecture logicielle en langage UML

![Première version](/architecture/projet_uml.png "Architecture du projet en UML")

## Conception

### Problématique

**_Apis Mellifera_ vs _Vespa crabo_**

Quelles proportions d’abeilles à miel faut-il à une colonie pour prospérer face au danger des frelons ?

Ouverture avec _Vespa velutina_ (frelons asiatiques)

### Courte description

Pour répondre à cette problématique, il nous faut simuler des abeilles à miel _Apis Mellifera_ et des frelons européens _Vespa crabo_. Les abeilles à miel vivent en colonies installées dans une ruche.
L'_Api Mellifera_ a un cycle de vie de 40 jours. Durant sa vie, elle change de métiers plusieurs fois, dans l'ordre : (nettoyeuse, nourice, architecte, faiseuse de miel, gardienne => METIERS DE RUCHE), butineuse.
L'abeille vient au monde adulte, avant cela elle est à l'état de larve dans une alvéole.
Les 28 premiers jours de sa vie, elle a un métier de ruche et passe les 12 derniers jours de sa vie en tant que butineuse (seul métier lui permettant de sortir de la ruche).
Il faut une certaine quantité de miel par jour à la ruche pour survivre.
La reine reste toujours dans la ruche. Une butineuse a une capacité de nectar et de pollen maximum qu'elle peut transporter jusqu'à la ruche.
Une butineuse rentre à la ruche à partir d'un certain niveau de nectar.
La production de miel nécessite une certaine quantité de nectar et un certain nombre d'abeilles d'intérieur.

### Règles

- Une abeille vit dans une ruche
- Une abeille vit 40 jours
- Une abeille a trois métiers au cours de sa vie : ouvrière, gardienne et butineuse
- Une abeille est ouvrière ou gardienne pendant 28 jours
- Une abeille est butineuse jusqu'à 12 jours à la fin de sa vie, seul métier lui permettant de sortir de la ruche
- Une butineuse trouve des fleurs à butiner
- Une butineuse avec sa capacité de nectar au maximum rentre déposer le nectar et le pollen à la ruche
- Une butineuse qui rencontre une fleur retient sa position
- Une butineuse connaissant la position d’au moins une fleur partage ses connaissances aux autres butineuses
- Une butineuse peut butiner une fleur à la fois et une fleur ne peut être butinée que par une butineuse à la fois
- Une gardienne reste à la ruche et ne produit pas de miel
- Une gardienne alerte toutes les abeilles de la ruche lors d'une attaque de frelon
- Une ouvrière transforme le nectar en miel
- Les abeilles fuient les frelons (mais défendent la colonie au prix de leur vie)
- La ruche contient une reine qui produit des abeilles en consommant du miel
- La ruche meurt si elle n’a pas suffisamment de miel par jour.
- La ruche produit une quantité de miel proportionnelle à la quantité de miel amenée par les butineuses et le nombre d'abeilles d'intérieur.
- Une ruche peut tuer un frelon si il y a suffisamment d’abeilles à l'intérieur (trouver le ratio en ligne)
- Un frelon chasse les abeilles
- Un frelon cherche les ruches d’abeilles
- Suffisamment de frelons peuvent tuer la reine
- Un frelon peut environner des fleurs, suivre des abeilles ou environner la ruche

On pourrait proposer à l'utilisateur de gérer la proportion d'abeilles productrices et exploratrices en temps réel (chaque abeille ayant un coût en miel pour la ruche).

Améliorations possibles :

- La ruche peut être déplacée contre une grosse quantité de miel (proportionnel au nombre d’habitants)

## Sources

### Danse des abeilles

- [« The Waggle Dance of the Honeybee »](https://www.youtube.com/watch?v=bFDGPgXtK-U), Georgia Tech College of Computing, 03/02/2011 (visionnée le 28/11/2024).
- [« Bee Dance Language - the linguistics behind animal language »](https://www.youtube.com/watch?v=pb1lRI-YePU), NativLang, 20/12/2014 (visionnée le 28/11/2024)
- [« Les moyens de communication impressionnants des abeilles »](https://www.nationalgeographic.fr/video/animaux/les-moyens-de-communication-impressionnants-des-abeilles), National Geographic, 14/09/2023 (visionnée le 28/11/2024)

### Généralités et apiculture

- Podcast Le bruit de la ruche (Saison 1, épisodes 1 à 5 écoutés le 28/11/2024)
- Documentaire « Des abeilles et des hommes - Et si les abeilles disparaissaient », Markus Imhoof, 2012 (visionné le 28/11/2024)
- http://rucherecole68.thann.free.fr/Calendrier%20apicole/des_abeilles_des_chiffres.pdf

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

### Notes

#### _Apis Mellifera_

Objectif des abeilles ouvrières : protéger la colonie.
Objectif de la reine : pondre des oeufs. (abstraire la reine en objet pour l'instant)
Objectif du faux bourdon : féconder la reine. (il meurt après, s'il ne féconde pas avant l'hiver il est tué car ne sert plus à rien)

Se décline en 3 types :

- Ouvrières (40 jours de vie) => oeuf fécondé
- Reine (jusqu'à 7 ans) : seul métier reine => oeuf fécondé
- Faux bourdon (jusqu'à 3 mois) => oeufs non fécondés

Métiers des abeilles :

- Nettoyeuse INTERIEUR
- Nourice (nourrit larves et produit Gelée royale) INTERIEUR
- Dame d'honneur (10-15) (nourrir la reine à la gelée royale, nettoyer la reine) INTERIEUR
- Ventileuse (maintenir température dans la ruche et assécher le miel) NON
- Architectes/cirières (au 12e jour) (produit de la cire à partir de leur corps et construire la ruche, les alvéoles) INTERIEUR
- Intendante (organise le rangement dans la ruche, miel au bon endroit, pollen au bon endroit, ponte de la reine) INTERIEUR
- Faiseuses de miel (reçoit nectar des butineuses et se transmettent le nectar de bouche à bouche pour créer le miel) INTERIEUR
- Gardienne (garde la ruche, le miel contre les papillons, bourdons, frelons, guèpes, ours, humains, souris, oiseau) INTERIEUR
- Butineuse (sort de la ruche, dernier métier, récupère nectar et pollen) EXTERIEUR
- Éclaireuses (les meilleures butineuses) EXTERIEUR => butineuse

Il faut 25kg de miel à une colonie pour passer un hiver normal

Une abeille ramène à la ruche nectar et pollen. C'est le nectar qui sert à faire le miel.

Quelle capacité de pollen par abeille ?
Nombre d'alvéoles dans une ruche ?

Passer la fleur en agent (Percept de la place dispo dans un périmètre pour se reproduire)

https://www.miel-direct.fr/bilan-annuel-colonie/ :

- butineuse peut prélever 40 mg de nectar, 25 sorties pour remplir une cellule
- Chaque abeille peut transporter une pelote de pollen sur chacune de ses pattes arrières (15 mg au total).
- la production de miel représente 60 kg environ par colonie et par an. Il en faut 40 kg pour couvrir les besoins internes en nourriture des abeilles. Selon les périodes de l’année, leurs besoins ne sont pas les mêmes. Lors de la saison de production de la gelée royale, les nourricières ont besoin de - 500 g de miel par jour. Les 20 kg supplémentaires sont du surplus que l’apiculteur peut se permettre de prélever.
- la reine pond entre 1000 et 2000 oeufs par jour
