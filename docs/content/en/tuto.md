---
title: Content styling tuto
linktitle: Content styling tuto
description:
date: 2018-09-05
publishdate: 2018-09-05
lastmod: 2018-09-05
categories: [content tuto]
keywords: []
menu:
  main:
    parent: "section name"
    weight: 01
weight: 3
draft: false
aliases: []
toc: true
---


# H1 Title
## H2 Title
### H3 Title
#### H4 Title


## Bullet Point

* Line 1
* Line 2
* Line 3
    * Line 31
    * line 32
* Line4


## Word with link

[EOS Canada](https://www.eoscanada.com)



## Words styling

*Text in italic*

**Text in bold**

***Text in bold and italic***



## Add note with vertical left color border

{{% note %}}
Our note here
{{% /note %}}



## Add note box with outside border


```bash
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec metus urna, euismod et ante ut, convallis viverra tellus. Aenean non tellus vel orci viverra lobortis a vel nunc. Curabitur luctus ac quam nec vulputate. Nam malesuada, diam a hendrerit iaculis, felis velit porttitor arcu, eu placerat sapien odio id nunc. Donec lorem elit, fermentum a faucibus ac, pulvinar id nunc. Integer commodo ex purus, non fermentum arcu viverra mollis. Phasellus faucibus ut urna non faucibus. Mauris sapien tortor, auctor et mi ut, consectetur scelerisque mauris. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos.
```



## Add note box with button to expend horizontal box

```
▶ hugo server -D

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec metus urna, euismod et ante ut, convallis viverra tellus. Aenean non tellus vel orci viverra lobortis a vel nunc. Curabitur luctus ac quam nec vulputate. Nam malesuada, diam a hendrerit iaculis, felis velit porttitor arcu, eu placerat sapien odio id nunc. Donec lorem elit, fermentum a faucibus ac, pulvinar id nunc. Integer commodo ex purus, non fermentum arcu viverra mollis. Phasellus faucibus ut urna non faucibus. Mauris sapien tortor, auctor et mi ut, consectetur scelerisque mauris. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos.
```




## Add video

{{< youtube id="FDy6RPAFMBA" >}}

{{< youtube id="aqeLuijgL3g" >}}

## add code with COPY button

{{< code file="archetype-example.sh" >}}
hugo new posts/my-first-post.md
{{< /code >}}



## add tweet post in page
{{< tweet 1005812880325922817 >}}


## Table

| Metric Name         | Description |
|---------------------|-------------|
| cumulative duration | The cumulative time spent executing a given template. |
| average duration    | The average time spent executing a given template. |
| maximum duration    | The maximum time a single execution took for a given template. |
| count               | The number of times a template was executed. |
| template            | The template name. |

