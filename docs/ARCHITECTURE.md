# GoPay Architecture

## Overview

GoPay is a payment orchestration system that provides a unified interface for multiple payment providers while maintaining its own state and data consistency.

## Current Architecture
<img width="1487" alt="Screenshot 2025-01-26 at 15 24 37" src="https://github.com/user-attachments/assets/bfe03b0b-af9c-4fc0-8720-c0373f5a913b" />


## Architecture Components

### 1. API Layer
- RESTful API endpoints for payments, subscriptions, disputes, and plans
- Authentication and authorization
- Rate limiting for API requests
- Request validation and error handling

### 2. Service Layer
- **Payment Service**: Handles payment processing and refunds
- **Subscription Service**: Manages recurring billing and subscription lifecycle
- **Dispute Service**: Handles payment disputes and evidence management
- **Configuration Service**: Manages system configuration and provider settings

### 3. Repository Layer
- Implements data access patterns using GORM
- Handles database operations and transactions
- Provides clean interfaces for services
- Manages relationships between entities

### 4. Provider Layer
- Abstract interface for payment providers
- Provider-specific implementations
- Handles provider API communication
- Maps provider responses to internal models

### 5. Data Layer
- PostgreSQL database with GORM as ORM
- Redis for caching and rate limiting
- Configuration storage
- Efficient querying and data retrieval

## Data Models

### Core Entities
1. **Payment**
   - Transaction details
   - Payment status tracking
   - Provider information
   - Refund history

2. **Subscription**
   - Recurring billing information
   - Subscription status
   - Plan details
   - Payment history

3. **Dispute**
   - Dispute details
   - Evidence management
   - Resolution tracking
   - Related transaction info

4. **Plan**
   - Pricing information
   - Billing intervals
   - Features and limits
   - Active status

## Implementation Details

### Database Access
- GORM for object-relational mapping
- Structured database schema
- Automated migrations
- Transaction support
- Relationship handling

### Data Flow
1. **API Request Flow**
   ```
   HTTP Request -> Handler -> Service -> Repository -> Database
                         -> Service -> Payment Provider
   ```

2. **Database Operations**
   ```
   Service Layer -> Repository Layer -> GORM -> PostgreSQL
   ```
