beverage = "tea"

switch (beverage) {
  case "water" {
    console.log("Water is $0.75 per bottle.")
  }
  case "juice" {
    console.log("Juice is $1.25 per bottle.")
  }
  case "coffee", "latte" {
    console.log("Coffee and lattes are $2.75 per 12oz.")
  }
  default {
    console.log("Unknown beverage.")
  }
}