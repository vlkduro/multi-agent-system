# Apis mellifera vs Vespa crabo

| Information | Value                                                                                  |
| ----------- | -------------------------------------------------------------------------------------- |
| Context     | UV AI30 (Multi-agent Systems) taught in Fall 2024 at UTC by Sylvain Lagrue             |
| Authors     | Quentin Valakou <br> Clément Martins <br> Julie Chartier <br> Alexandre Bidaux          |
| Date        | 11/27/24                                                                               |
| Version     | 1.0.0                                                                                  |

## Table of Contents

- [Project](#project)
- [Results](#results)
- [Installation](#installation)
- [Design](#design)
- [Sources](#sources)

## Project

### Problem Statement

**_Apis mellifera_ vs _Vespa crabo_**

What proportion of honey bees does a colony need to thrive in the face of hornet threats?

Our simulation observes whether a bee colony can survive thanks to foragers that collect nectar from flowers and workers that transform nectar into honey inside the hive. We run simulations to see whether the colony can survive against one or several hornets, depending on parameters such as: map size, number of hornets and bees at the start, number of flowers, maximum nectar capacity per bee, etc.

As an extension, we could integrate a more aggressive predator that honey bees are less able to defend against: the _Vespa velutina_ (Asian hornet).

### Rules

To answer this question, we simulate honey bees (_Apis mellifera_) and European hornets (_Vespa crabo_). Honey bees live in colonies housed in a hive.

- A bee lives in a hive.  
- A bee has three roles during its life: worker, guard, forager.  
- Only foragers can leave the hive.  
- A forager finds flowers to collect nectar.  
- A forager brings nectar back to the hive according to her carrying capacity.  
- A forager that finds a flower remembers its position.  
- A forager knowing at least one flower shares its position with other foragers.  
- A forager can collect from only one flower at a time, and each flower can only be collected by one forager at a time.  
- A guard stays at the hive and does not produce honey.  
- A guard alerts all bees in the hive when hornets attack.  
- A worker transforms nectar into honey.  
- Bees flee from hornets (but defend the colony at the cost of their lives).  
- The hive contains a queen that produces bees by consuming honey.  
- The hive dies if the queen is lost.  
- A hive can kill a hornet if enough bees are present.  
- A hornet hunts bees.  

The user plays the role of the queen, managing the production of bees in the hive (each bee costs honey to produce).

**Possible Improvements**:  
- The hive can be moved for a large honey cost (proportional to the number of inhabitants).

---

## Results

### Defense Presentation

[Presentation PDF](/doc/Apis%20Mellifera%20vs%20Vespa%20crabo%20Soutenance.pdf)

### Demo Screenshots

![Demo Screenshot](/doc/img/Screenshot%202025-01-06%20at%2023-26-07%20AI30%20WebSocket.png)  
![Demo Screenshot](/doc/img/Screenshot%202025-01-06%20at%2023-30-53%20AI30%20WebSocket.png)  

We can distinguish three main emergent behaviors in these screenshots and simulations:

- Forager bee trails emerge toward flower clusters, enabled by communication within the hive about flower positions.  
- Hornets converge toward bee trails.  
- Bees change objectives when encountering hornets, creating new trails toward other flowers. They no longer share disrupted paths with other foragers.  

---

## Installation

From a terminal, perform the following steps:

### Clone the project

With HTTPS:
```bash
git clone https://github.com/vlkduro/multi-agent-system.git
```

With SSH:
```bash
git clone git@github.com:vlkduro/multi-agent-system.git
```

### Launch the simulation from the project root
```bash
go run .\cmd\launch-server.go
```

### Open the interface in a browser
[http://localhost:8000/](http://localhost:8000/)

---

## Design

### UML Model

![](/doc/conception/simu_bees_uml_v3.png "UML Project Model")

---

## Sources

### Bee Dance

- [“The Waggle Dance of the Honeybee”](https://www.youtube.com/watch?v=bFDGPgXtK-U), Georgia Tech College of Computing, 02/03/2011 (viewed 11/28/2024).  
- [“Bee Dance Language - the linguistics behind animal language”](https://www.youtube.com/watch?v=pb1lRI-YePU), NativLang, 12/20/2014 (viewed 11/28/2024).  
- [“The Impressive Communication Methods of Bees”](https://www.nationalgeographic.fr/video/animaux/les-moyens-de-communication-impressionnants-des-abeilles), National Geographic, 09/14/2023 (viewed 11/28/2024).  
- [“Honey bees perfect their waggle dances by learning from elders”](https://www.youtube.com/watch?v=Cylim87fFgU), Science Magazine, 04/25/2023 (to be viewed).  

### General and Beekeeping

- Podcast *Le bruit de la ruche* (Season 1, episodes 1–5 listened to on 11/28/2024).  
- Documentary *More than Honey*, Markus Imhoof, 2012 (viewed 11/28/2024).  
- Frédéric Schubnel, Robert Hummel, Maurice Feltin, [“Des abeilles et des chiffres”](http://rucherecole68.thann.free.fr/Calendrier%20apicole/des_abeilles_des_chiffres.pdf), 09/2019 (consulted 12/02/2024).  
- [“The Different Job Titles of a Worker Bee”](https://bee2bee.com.au/the-different-job-titles-of-a-worker-bee/) (consulted 12/02/2024).  
- [“Annual Hive Figures: Colony Yearly Balance”](https://www.miel-direct.fr/bilan-annuel-colonie/) (consulted 12/02/2024).  
- [“‘Bees are sentient’: inside the stunning brains of nature’s hardest workers”](https://www.theguardian.com/environment/2023/apr/02/bees-intelligence-minds-pollination) (consulted 11/28/2024).  
- [“‘Like sending bees to war’: the deadly truth behind your almond milk obsession”](https://www.theguardian.com/environment/2020/jan/07/honeybees-deaths-almonds-hives-aoe) (consulted 11/28/2024).  
- [“Apis mellifera – European honey bee”](https://fr.wikipedia.org/wiki/Apis_mellifera) (consulted 11/29/2024).  

### Bee Predators

- [“Predators”](https://www.apiculteurs-occitanie.fr/informations/predateurs/) (consulted 11/28/2024).  
- [“Colony of invasive red dwarf honeybee found for first time in Europe”](https://www.theguardian.com/environment/article/2024/aug/30/colony-invasive-red-dwarf-honeybee-first-time-europe) (consulted 11/28/2024).  
- [“Asian hornets overwintered in UK for first time, DNA testing shows”](https://www.theguardian.com/environment/article/2024/jun/03/asian-hornets-overwintered-uk-first-time-dna-testing) (consulted 11/28/2024).  
- [“How to distinguish Asian hornets from European hornets”](https://www.apiculture.net/blog/savoir-differencier-le-frelon-asiatique-de-l-europeen-n181) (consulted 11/29/2024).  
- [“Better understanding the European hornet”](https://www.apiculture.net/blog/frelon-europeen-n104) (consulted 11/29/2024).  
- [“The Asian hornet: a serious threat to our local bee colonies”](https://www.pollinis.org/publications/le-frelon-asiatique-un-fleau-redoutable-pour-nos-colonies-d-abeilles-locales/) (consulted 11/29/2024).  
- [“Neature – European hornet”](https://neature.fr/especes/frelon-europeen/) (consulted 12/02/2024).  
- [“Vespa (insect) – Hornet”](https://fr.wikipedia.org/wiki/Vespa_(insecte)#Ph%C3%A9romone_de_reconnaissance_et_de_d%C3%A9fense_group%C3%A9e) (consulted 11/29/2024).  

### Defense Against Predators

- [“The impressive defense system of the Japanese honeybee against a hornet”](https://www.youtube.com/watch?v=bDHWv8_-f4U), National Geographic Wild France, 11/21/2020 (viewed 11/28/2024).  
- [“This hornet shouldn’t have messed with these bees”](https://www.youtube.com/watch?v=N_AO8XvAZgQ), National Geographic Wild France, 04/28/2021 (viewed 11/28/2024).  
- [“How do bees manage to defend themselves?”](https://www.nationalgeographic.fr/video/animaux/comment-les-abeilles-parviennent-elles-a-se-defendre), National Geographic, 04/15/2023 (viewed 11/28/2024).  

### Others

- [“How Far do Honeybees Fly?”](https://wildflowermeadows.com/2024/01/how-far-do-honeybees-fly/)  
- [“The Sunflowers”](https://www.miel-lerucherdelours.fr/en/content/67-the-sunflowers)  
- [“Pathfinding for Beginners”](https://web.archive.org/web/20171022224528/http://www.policyalmanac.org:80/games/aStarTutorial.htm)  
- [“How to Find the Intersection Points of Two Circles with Different Centres and Radiuses”](https://www.youtube.com/watch?v=PSlWb90JJx4)  
- [“Check whether a given point lies inside a triangle or not”](https://www.geeksforgeeks.org/check-whether-a-given-point-lies-inside-a-triangle-or-not/)  
- [“What is the correct way to find the min between two integers in Go?”](https://stackoverflow.com/questions/27516387/what-is-the-correct-way-to-find-the-min-between-two-integers-in-go)  
- [“Round numbers in Go”](https://dev.to/natamacm/round-numbers-in-go-5c01)  
