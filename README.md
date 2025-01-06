# Apis mellifera vs Vespa crabo

| Information | Valeur                                                                         |
| ----------- | ------------------------------------------------------------------------------ |
| Contexte    | UV AI30 (Systèmes multi-agents) dispensée en A24 à l'UTC par Sylvain Lagrue    |
| Auteurs     | Valakou Quentin <br> Martins Clément <br> Chartier Julie <br> Bidaux Alexandre |
| Date        | 27/11/24                                                                       |
| Version     | 0.0.1                                                                          |

## Sommaire

- [Résultat](#résultat)
- [Installation](#installation)
- [Conception](#conception)
- [Architecture](#architecture)
- [Sources](#sources)

## Résultat

### Lien vers la présentation de la soutenance

[Google slides](https://docs.google.com/presentation/d/1iAHNwlmZweJUNYqqZqQWgXstjSIhNtLi2vJ8IsF_iug/edit?usp=sharing)

### Capture d'écran de la démo

![Capture d'écran de la démo](/docs/)

## Installation

Depuis un terminal effectuer les étapes suivantes.

### Cloner le projet

Avec HTTPS :

```
git clone https://gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux.git
```

ou avec SSH :

```
git clone git@gitlab.utc.fr:bidauxal/ai30_valakou_martins_chartier_bidaux.git
```

### Lancer la simulation depuis la racine du projet

```
go run .\cmd\launch-server.go
```

## Conception

### Problématique

**_Apis Mellifera_ vs _Vespa crabo_**

Quelles proportions d’abeilles à miel faut-il à une colonie pour prospérer face au danger des frelons ?

Ouverture avec _Vespa velutina_ (frelons asiatiques)

### Courte description

Pour répondre à cette problématique, il nous faut simuler des abeilles à miel _Apis Mellifera_ et des frelons européens _Vespa crabo_. Les abeilles à miel vivent en colonies installées dans une ruche.

### Règles

- Une abeille vit dans une ruche.
- Une abeille a trois métiers au cours de sa vie : ouvrière, gardienne, butineuse.
- Une butineuse est la seule pouvant sortir de la ruche.
- Une butineuse trouve des fleurs à butiner.
- Une butineuse dépose le nectar à la ruche en fonction de la quantité qu’elle peut transporter.
- Une butineuse qui rencontre une fleur retient sa position.
- Une butineuse connaissant la position d’au moins une fleur partage ses connaissances aux autres butineuses.
- Une butineuse peut butiner une fleur à la fois et une fleur ne peut être butinée que par une butineuse à la fois.
- Une gardienne reste à la ruche et ne produit pas de miel.
- Une gardienne alerte toutes les abeilles de la ruche lors d'une attaque de frelon.
- Une ouvrière transforme le nectar en miel.
- Les abeilles fuient les frelons (mais défendent la colonie au prix de leur vie).
- La ruche contient une reine qui produit des abeilles en consommant du miel.
- La ruche meurt si elle n’a plus de reine.
- Une ruche peut tuer un frelon si il y a suffisamment d’abeilles à l'intérieur.
- Un frelon chasse les abeilles.
- Un frelon cherche les ruches d’abeilles.
- Suffisamment de frelons peuvent tuer la reine.
- Un frelon ayant tué suffisamment d'abeilles se dirigera vers la ruche.
- Un frelon peut environner des fleurs, suivre des abeilles ou environner la ruche.

On pourrait proposer à l'utilisateur de gérer la proportion d'abeilles productrices et exploratrices en temps réel (chaque abeille ayant un coût en miel pour la ruche).

Améliorations possibles :

- La ruche peut être déplacée contre une grosse quantité de miel (proportionnel au nombre d’habitants)

### Modélisation UML

![](/docs/conception/simu_bees_uml_v2.png "Modélisation du projet en UML")

## Architecture

`Agent` et `Object` sont des interfaces de même stabilité que l'environement. C'est à dire que lorsque une des interfaces change, `Environment` change également. Ils sont donc dans le même package pour l'instant

> Remarque : mettre `Agent` (i.e `Object`) dans le package `agent` (i.e `object`) crée une boucle de dépendance :
> `Agent` appelle le package `environment` et `Environment` appelle le package `agent`, par défaut le compilateur interdit les dépendances cycliques

## Sources

### Danse des abeilles

- [« The Waggle Dance of the Honeybee »](https://www.youtube.com/watch?v=bFDGPgXtK-U), Georgia Tech College of Computing, 03/02/2011 (visionnée le 28/11/2024).
- [« Bee Dance Language - the linguistics behind animal language »](https://www.youtube.com/watch?v=pb1lRI-YePU), NativLang, 20/12/2014 (visionnée le 28/11/2024)
- [« Les moyens de communication impressionnants des abeilles »](https://www.nationalgeographic.fr/video/animaux/les-moyens-de-communication-impressionnants-des-abeilles), National Geographic, 14/09/2023 (visionnée le 28/11/2024)
- [« Honey bees perfect their waggle dances by learning from elders »](https://www.youtube.com/watch?v=Cylim87fFgU), Science Magazine, 25/04/2023 (à visionner).

### Généralités et apiculture

- Podcast Le bruit de la ruche (Saison 1, épisodes 1 à 5 écoutés le 28/11/2024)
- Documentaire « Des abeilles et des hommes - Et si les abeilles disparaissaient », Markus Imhoof, 2012 (visionné le 28/11/2024)
- Frédéric Schubnel, Robert Hummel, Maurice Feltin, [« Des abeilles et des chiffres »](http://rucherecole68.thann.free.fr/Calendrier%20apicole/des_abeilles_des_chiffres.pdf), Syndicat des apiculteurs de Thann et environs, 09/2019 (consulté le 02/12/2024).
- https://bee2bee.com.au/the-different-job-titles-of-a-worker-bee/ (consulté le 02/12/2024)
- https://www.miel-direct.fr/bilan-annuel-colonie/ (consulté le 02/12/2024)
- https://www.theguardian.com/environment/2023/apr/02/bees-intelligence-minds-pollination (consulté le 28/11/2024)
- https://www.theguardian.com/environment/2020/jan/07/honeybees-deaths-almonds-hives-aoe (consulté le 28/11/2024)
- https://fr.wikipedia.org/wiki/Apis_mellifera (consulté le 29/11/2024)

### Prédateurs des abeilles

- https://www.apiculteurs-occitanie.fr/informations/predateurs/ (consulté le 28/11/2024) : prédateurs des abeilles
- https://www.theguardian.com/environment/article/2024/aug/30/colony-invasive-red-dwarf-honeybee-first-time-europe (consulté le 28/11/2024)
- https://www.theguardian.com/environment/article/2024/jun/03/asian-hornets-overwintered-uk-first-time-dna-testing (consulté le 28/11/2024)
- https://www.apiculture.net/blog/savoir-differencier-le-frelon-asiatique-de-l-europeen-n181 (consulté le 29/11/2024)
- https://www.apiculture.net/blog/frelon-europeen-n104 (consulté le 29/11/2024)
- https://www.pollinis.org/publications/le-frelon-asiatique-un-fleau-redoutable-pour-nos-colonies-d-abeilles-locales/ (consulté le 29/11/2024)
- https://neature.fr/especes/frelon-europeen/ (consulté le 02/12/2024)
- https://fr.wikipedia.org/wiki/Vespa_(insecte)#Ph%C3%A9romone_de_reconnaissance_et_de_d%C3%A9fense_group%C3%A9e (consulté le 29/11/2024)

### Défense contre les prédateurs

- [« Le redoutable système de défense de l'abeille japonaise face à un frelon »](https://www.youtube.com/watch?v=bDHWv8_-f4U), National Geographic Wild France, 21/11/2020 (visionnée le 28/11/2024)
- [« Ce frelon n'aurait pas dû se frotter à ces abeilles »](https://www.youtube.com/watch?v=N_AO8XvAZgQ), National Geographic Wild France, 28/04/2021 (visionnée le 28/11/2024)
- [« Comment les abeilles parviennent-elles à se défendre ? »](https://www.nationalgeographic.fr/video/animaux/comment-les-abeilles-parviennent-elles-a-se-defendre), National Geographic, 15/04/2023 (visionnée le 28/11/2024)

### Implémentation vision et pathfinding

-
