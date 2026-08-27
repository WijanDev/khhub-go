# Cost analysis

- **Score:** 3
- **Scale:** 1 = cheap … 5 = expensive
- **Ratio (utility ÷ cost):** 1.33

## Justification

Steps 1–2 are org/project setup and a properties file. Step 3 changes CI (coverage artifacts + scan). Step 4 is a short doc. Step 5 is the gate in the Sonar UI. Several days including token, first noisy scan, and exclusion tuning. Not a product schema change. Risk is a gate that fails every PR until exclusions are right.
