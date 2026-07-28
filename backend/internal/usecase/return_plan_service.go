package usecase

import (
	"context"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IReturnPlanService interface {
	CreateDraft(ctx context.Context, userID uuid.UUID, targetPet, targetCan, targetGlass int) (entity.ReturnPlanWithPoints, error)
	GetByID(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlanWithPoints, error)
	GetDraft(ctx context.Context, userID uuid.UUID) (entity.ReturnPlanWithPoints, error)
	List(ctx context.Context, userID uuid.UUID) ([]entity.ReturnPlan, error)
	UpdateTargets(ctx context.Context, userID, planID uuid.UUID, targetPet, targetCan, targetGlass int) (entity.ReturnPlanWithPoints, error)
	Save(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlanWithPoints, error)
	Delete(ctx context.Context, userID, planID uuid.UUID) error
	AddPoint(ctx context.Context, userID, planID, returnPointID uuid.UUID, petCount, canCount, glassCount int) (entity.ReturnPlanPoint, error)
	UpdatePointQuantities(ctx context.Context, userID, planID, planPointID uuid.UUID, petCount, canCount, glassCount int) (entity.ReturnPlanPoint, error)
	DeletePoint(ctx context.Context, userID, planID, planPointID uuid.UUID) error
	ConfirmPoint(ctx context.Context, userID, planID, planPointID uuid.UUID) (entity.ReturnPlanWithPoints, error)
	ConfirmAll(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlanWithPoints, error)
}

type returnPlanService struct {
	repo            IReturnPlanRepository
	returnPointRepo IReturnPointRepository
}

func NewReturnPlanService(repo IReturnPlanRepository, returnPointRepo IReturnPointRepository) IReturnPlanService {
	return &returnPlanService{
		repo:            repo,
		returnPointRepo: returnPointRepo,
	}
}

func (s *returnPlanService) CreateDraft(ctx context.Context, userID uuid.UUID, targetPet, targetCan, targetGlass int) (entity.ReturnPlanWithPoints, error) {
	if userID == uuid.Nil {
		return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("user id is required")
	}
	if err := validateTargets(targetPet, targetCan, targetGlass); err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	now := time.Now().UTC()
	plan, err := s.repo.CreateReturnPlan(ctx, entity.ReturnPlan{
		ID:          uuid.New(),
		UserID:      &userID,
		Status:      entity.PlanStatusDraft,
		TargetPet:   targetPet,
		TargetCan:   targetCan,
		TargetGlass: targetGlass,
		RatePet:     entity.DefaultRatePetGrosze,
		RateCan:     entity.DefaultRateCanGrosze,
		RateGlass:   entity.DefaultRateGlassGrosze,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	return entity.ReturnPlanWithPoints{Plan: plan, Points: []entity.ReturnPlanPoint{}}, nil
}

func (s *returnPlanService) GetByID(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlanWithPoints, error) {
	return s.getOwnedPlanWithPoints(ctx, userID, planID)
}

func (s *returnPlanService) GetDraft(ctx context.Context, userID uuid.UUID) (entity.ReturnPlanWithPoints, error) {
	if userID == uuid.Nil {
		return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("user id is required")
	}

	plan, err := s.repo.GetDraftReturnPlanByUserID(ctx, userID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	points, err := s.repo.ListReturnPlanPointsByPlanID(ctx, plan.ID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	return entity.ReturnPlanWithPoints{Plan: plan, Points: points}, nil
}

func (s *returnPlanService) List(ctx context.Context, userID uuid.UUID) ([]entity.ReturnPlan, error) {
	if userID == uuid.Nil {
		return []entity.ReturnPlan{}, entity.NewBadRequest("user id is required")
	}

	plans, err := s.repo.ListReturnPlansByUserID(ctx, userID)
	if err != nil {
		return []entity.ReturnPlan{}, err
	}

	return plans, nil
}

func (s *returnPlanService) UpdateTargets(ctx context.Context, userID, planID uuid.UUID, targetPet, targetCan, targetGlass int) (entity.ReturnPlanWithPoints, error) {
	if err := validateTargets(targetPet, targetCan, targetGlass); err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	plan, err := s.requireOwnedEditablePlan(ctx, userID, planID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	points, err := s.repo.ListReturnPlanPointsByPlanID(ctx, plan.ID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	assigned := sumPointQuantities(points)
	if assigned.pet > targetPet || assigned.can > targetCan || assigned.glass > targetGlass {
		return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("targets cannot be lower than already assigned quantities")
	}

	updated, err := s.repo.UpdateReturnPlanTargets(ctx, plan.ID, targetPet, targetCan, targetGlass)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	return entity.ReturnPlanWithPoints{Plan: updated, Points: points}, nil
}

func (s *returnPlanService) Save(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlanWithPoints, error) {
	plan, err := s.requireOwnedPlan(ctx, userID, planID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	if plan.Status == entity.PlanStatusConfirmed {
		return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("confirmed plan cannot be saved again")
	}

	if plan.Status == entity.PlanStatusDraft {
		plan, err = s.repo.UpdateReturnPlanStatus(ctx, plan.ID, entity.PlanStatusSaved)
		if err != nil {
			return entity.ReturnPlanWithPoints{}, err
		}
	}

	points, err := s.repo.ListReturnPlanPointsByPlanID(ctx, plan.ID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	return entity.ReturnPlanWithPoints{Plan: plan, Points: points}, nil
}

func (s *returnPlanService) Delete(ctx context.Context, userID, planID uuid.UUID) error {
	if userID == uuid.Nil || planID == uuid.Nil {
		return entity.NewBadRequest("user id and plan id are required")
	}

	plan, err := s.requireOwnedPlan(ctx, userID, planID)
	if err != nil {
		return err
	}

	if plan.Status == entity.PlanStatusConfirmed {
		return entity.NewBadRequest("confirmed plan cannot be deleted")
	}

	return s.repo.DeleteReturnPlanByIDAndUserID(ctx, planID, userID)
}

func (s *returnPlanService) AddPoint(ctx context.Context, userID, planID, returnPointID uuid.UUID, petCount, canCount, glassCount int) (entity.ReturnPlanPoint, error) {
	if returnPointID == uuid.Nil {
		return entity.ReturnPlanPoint{}, entity.NewBadRequest("return point id is required")
	}
	if err := validateQuantities(petCount, canCount, glassCount); err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	plan, err := s.requireOwnedEditablePlan(ctx, userID, planID)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	returnPoint, err := s.returnPointRepo.GetByID(ctx, returnPointID)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	if err := validateAcceptedTypes(returnPoint, petCount, canCount, glassCount); err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	points, err := s.repo.ListReturnPlanPointsByPlanID(ctx, plan.ID)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	if err := validateAllocationFits(plan, points, uuid.Nil, petCount, canCount, glassCount); err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	now := time.Now().UTC()
	point, err := s.repo.CreateReturnPlanPoint(ctx, entity.ReturnPlanPoint{
		ID:            uuid.New(),
		ReturnPlanID:  plan.ID,
		ReturnPointID: returnPoint.ID,
		PointName:     returnPoint.Name,
		PointAddress:  returnPoint.Address,
		PetCount:      petCount,
		CanCount:      canCount,
		GlassCount:    glassCount,
		Status:        entity.PlanPointStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	return point, nil
}

func (s *returnPlanService) UpdatePointQuantities(ctx context.Context, userID, planID, planPointID uuid.UUID, petCount, canCount, glassCount int) (entity.ReturnPlanPoint, error) {
	if err := validateQuantities(petCount, canCount, glassCount); err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	plan, err := s.requireOwnedEditablePlan(ctx, userID, planID)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	point, err := s.requirePlanPoint(ctx, plan.ID, planPointID)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}
	if point.Status != entity.PlanPointStatusPending {
		return entity.ReturnPlanPoint{}, entity.NewBadRequest("confirmed plan point cannot be updated")
	}

	returnPoint, err := s.returnPointRepo.GetByID(ctx, point.ReturnPointID)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}
	if err := validateAcceptedTypes(returnPoint, petCount, canCount, glassCount); err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	points, err := s.repo.ListReturnPlanPointsByPlanID(ctx, plan.ID)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}
	if err := validateAllocationFits(plan, points, point.ID, petCount, canCount, glassCount); err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	updated, err := s.repo.UpdateReturnPlanPointQuantities(ctx, point.ID, petCount, canCount, glassCount)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}

	return updated, nil
}

func (s *returnPlanService) DeletePoint(ctx context.Context, userID, planID, planPointID uuid.UUID) error {
	plan, err := s.requireOwnedEditablePlan(ctx, userID, planID)
	if err != nil {
		return err
	}

	point, err := s.requirePlanPoint(ctx, plan.ID, planPointID)
	if err != nil {
		return err
	}
	if point.Status != entity.PlanPointStatusPending {
		return entity.NewBadRequest("confirmed plan point cannot be deleted")
	}

	if err := s.repo.DeleteReturnPlanPoint(ctx, point.ID); err != nil {
		return err
	}

	return s.refreshPlanStatus(ctx, plan.ID)
}

func (s *returnPlanService) ConfirmPoint(ctx context.Context, userID, planID, planPointID uuid.UUID) (entity.ReturnPlanWithPoints, error) {
	plan, err := s.requireOwnedEditablePlan(ctx, userID, planID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	point, err := s.requirePlanPoint(ctx, plan.ID, planPointID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}
	if point.Status != entity.PlanPointStatusPending {
		return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("plan point is already confirmed")
	}
	if point.PetCount+point.CanCount+point.GlassCount == 0 {
		return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("cannot confirm plan point with zero quantities")
	}

	if _, err := s.repo.ConfirmReturnPlanPoint(ctx, point.ID); err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	if err := s.refreshPlanStatus(ctx, plan.ID); err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	return s.getOwnedPlanWithPoints(ctx, userID, plan.ID)
}

func (s *returnPlanService) ConfirmAll(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlanWithPoints, error) {
	plan, err := s.requireOwnedEditablePlan(ctx, userID, planID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	points, err := s.repo.ListReturnPlanPointsByPlanID(ctx, plan.ID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}
	if len(points) == 0 {
		return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("plan has no points to confirm")
	}

	assigned := sumPointQuantities(points)
	if assigned.pet != plan.TargetPet || assigned.can != plan.TargetCan || assigned.glass != plan.TargetGlass {
		return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("plan is incomplete: assigned quantities must match targets")
	}

	for _, point := range points {
		if point.Status == entity.PlanPointStatusPending && point.PetCount+point.CanCount+point.GlassCount == 0 {
			return entity.ReturnPlanWithPoints{}, entity.NewBadRequest("cannot confirm plan point with zero quantities")
		}
	}

	if _, err := s.repo.ConfirmAllPendingReturnPlanPoints(ctx, plan.ID); err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	if _, err := s.repo.ConfirmReturnPlan(ctx, plan.ID); err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	return s.getOwnedPlanWithPoints(ctx, userID, plan.ID)
}

func (s *returnPlanService) getOwnedPlanWithPoints(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlanWithPoints, error) {
	plan, err := s.requireOwnedPlan(ctx, userID, planID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	points, err := s.repo.ListReturnPlanPointsByPlanID(ctx, plan.ID)
	if err != nil {
		return entity.ReturnPlanWithPoints{}, err
	}

	return entity.ReturnPlanWithPoints{Plan: plan, Points: points}, nil
}

func (s *returnPlanService) requireOwnedPlan(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlan, error) {
	if userID == uuid.Nil || planID == uuid.Nil {
		return entity.ReturnPlan{}, entity.NewBadRequest("user id and plan id are required")
	}

	plan, err := s.repo.GetReturnPlanByIDAndUserID(ctx, planID, userID)
	if err != nil {
		return entity.ReturnPlan{}, err
	}

	return plan, nil
}

func (s *returnPlanService) requireOwnedEditablePlan(ctx context.Context, userID, planID uuid.UUID) (entity.ReturnPlan, error) {
	plan, err := s.requireOwnedPlan(ctx, userID, planID)
	if err != nil {
		return entity.ReturnPlan{}, err
	}
	if plan.Status == entity.PlanStatusConfirmed {
		return entity.ReturnPlan{}, entity.NewBadRequest("confirmed plan cannot be modified")
	}

	return plan, nil
}

func (s *returnPlanService) requirePlanPoint(ctx context.Context, planID, planPointID uuid.UUID) (entity.ReturnPlanPoint, error) {
	if planPointID == uuid.Nil {
		return entity.ReturnPlanPoint{}, entity.NewBadRequest("plan point id is required")
	}

	point, err := s.repo.GetReturnPlanPointByID(ctx, planPointID)
	if err != nil {
		return entity.ReturnPlanPoint{}, err
	}
	if point.ReturnPlanID != planID {
		return entity.ReturnPlanPoint{}, entity.NewNotFound("resource not found")
	}

	return point, nil
}

func (s *returnPlanService) refreshPlanStatus(ctx context.Context, planID uuid.UUID) error {
	plan, err := s.repo.GetReturnPlanByID(ctx, planID)
	if err != nil {
		return err
	}
	if plan.Status == entity.PlanStatusConfirmed {
		return nil
	}

	points, err := s.repo.ListReturnPlanPointsByPlanID(ctx, planID)
	if err != nil {
		return err
	}

	confirmedCount := 0
	for _, point := range points {
		if point.Status == entity.PlanPointStatusConfirmed {
			confirmedCount++
		}
	}

	var next entity.PlanStatus
	switch {
	case len(points) > 0 && confirmedCount == len(points):
		_, err = s.repo.ConfirmReturnPlan(ctx, planID)
		return err
	case confirmedCount > 0:
		next = entity.PlanStatusPartiallyConfirmed
	case plan.Status == entity.PlanStatusDraft:
		return nil
	default:
		next = entity.PlanStatusSaved
	}

	if plan.Status == next {
		return nil
	}

	_, err = s.repo.UpdateReturnPlanStatus(ctx, planID, next)
	return err
}

type quantitySum struct {
	pet   int
	can   int
	glass int
}

func sumPointQuantities(points []entity.ReturnPlanPoint) quantitySum {
	var sum quantitySum
	for _, point := range points {
		sum.pet += point.PetCount
		sum.can += point.CanCount
		sum.glass += point.GlassCount
	}
	return sum
}

func validateTargets(pet, can, glass int) error {
	if pet < 0 || can < 0 || glass < 0 {
		return entity.NewBadRequest("targets cannot be negative")
	}
	return nil
}

func validateQuantities(pet, can, glass int) error {
	if pet < 0 || can < 0 || glass < 0 {
		return entity.NewBadRequest("quantities cannot be negative")
	}
	return nil
}

func validateAcceptedTypes(point entity.ReturnPoint, pet, can, glass int) error {
	if pet > 0 && !point.Pet {
		return entity.NewBadRequest("return point does not accept PET")
	}
	if can > 0 && !point.Can {
		return entity.NewBadRequest("return point does not accept cans")
	}
	if glass > 0 && !point.Glass {
		return entity.NewBadRequest("return point does not accept glass")
	}
	return nil
}

func validateAllocationFits(plan entity.ReturnPlan, points []entity.ReturnPlanPoint, excludePointID uuid.UUID, pet, can, glass int) error {
	assigned := quantitySum{}
	for _, point := range points {
		if point.ID == excludePointID {
			continue
		}
		assigned.pet += point.PetCount
		assigned.can += point.CanCount
		assigned.glass += point.GlassCount
	}

	if assigned.pet+pet > plan.TargetPet {
		return entity.NewBadRequest("PET allocation exceeds plan target")
	}
	if assigned.can+can > plan.TargetCan {
		return entity.NewBadRequest("can allocation exceeds plan target")
	}
	if assigned.glass+glass > plan.TargetGlass {
		return entity.NewBadRequest("glass allocation exceeds plan target")
	}

	return nil
}
