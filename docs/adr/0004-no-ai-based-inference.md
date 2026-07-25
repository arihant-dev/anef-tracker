# ADR-0004: No AI-Based Inference

## Status
Accepted

## Context
Predictive ML models and AI inference introduce probabilistic uncertainty into administrative status tracking, which can mislead users regarding prefetoral decisions.

## Decision
We mandate **Deterministic Workflow Intelligence Only**:
1. All workflow state machine transitions, medians, and duration statistics are strictly derived from observed database facts.
2. No ML models, probabilistic inference, or hallucinations are used.
3. Every claim in generated reports must answer: "Which stored artifact proves this?"

## Consequences
- 100% explainable, reproducible, and verifiable operational intelligence.
- Eliminates false hope or inaccurate administrative progression predictions.
