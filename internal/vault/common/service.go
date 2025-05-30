package vault

import (
	"context"
	"errors"
	"fmt"
	"swagtask/internal/auth"
	db "swagtask/internal/db/generated"
	"swagtask/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---- READ ----

func GetVaultWithCollaboratorsById(queries *db.Queries, userId, id pgtype.UUID, ctx context.Context) (*VaultWithCollaborators, error) {
	vaultWithRelations, err := queries.GetVaultWithCollaborators(ctx, db.GetVaultWithCollaboratorsParams{
		UserID:  userId,
		VaultID: id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}
	if len(vaultWithRelations) == 0 {
		return nil, utils.ErrNotFound
	}
	var vault VaultUI
	collaboratorsOfVault := []RelatedCollaborator{}
	var pathToPfp string
	var username string
	for _, vaultWithRelation := range vaultWithRelations {
		if vaultWithRelation.CollaboratorUsername.Valid && vaultWithRelation.CollaboratorPathToPfp.Valid {
			collaboratorsOfVault = append(collaboratorsOfVault,
				RelatedCollaborator{
					Name:      vaultWithRelation.CollaboratorUsername.String,
					PathToPfp: vaultWithRelation.CollaboratorPathToPfp.String,
					Role:      string(vaultWithRelation.CollaboratorRole.VaultRelRoleType),
				})
		}

		if vaultWithRelation.CollaboratorRole.Valid && vaultWithRelation.CollaboratorRole.VaultRelRoleType == "owner" {
			pathToPfp = vaultWithRelation.CollaboratorPathToPfp.String
			username = vaultWithRelation.CollaboratorUsername.String
		}

		vault = NewUIvault(
			db.Vault{
				ID:          vaultWithRelation.ID,
				Name:        vaultWithRelation.Name,
				Description: vaultWithRelation.Description,
				Locked:      vaultWithRelation.Locked,
				UpdatedAt:   vaultWithRelation.UpdatedAt,
				CreatedAt:   vaultWithRelation.CreatedAt,
				Kind:        vaultWithRelation.Kind,
			},
			auth.Author{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		)

	}

	vaultWithCollaborators := newVaultWithCollaborators(vault, collaboratorsOfVault)
	return &vaultWithCollaborators, nil
}

func GetVaultsWithCollaborators(queries *db.Queries,
	userId pgtype.UUID, ctx context.Context) ([]VaultWithCollaborators, error) {

	vaultsWithCollaborators, err := queries.GetVaultsWithCollaborators(ctx, userId)
	if err != nil {
		fmt.Println("error here at first")
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	// maybe reimplement when adding vaults is done
	// allCollaborators, errAllCollaborators := queries.GetAllCollaboratorsDesc(ctx, db.GetAllCollaboratorsDescParams{
	// 	VaultID: vaultId,
	// 	UserID:  userId,
	// })
	// if errAllCollaborators != nil {
	// 	fmt.Println("error here at second")
	// 	return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errAllCollaborators)
	// }

	vaultsWithCollaboratorsUI := []VaultWithCollaborators{}
	vaultIdToCollaborators := make(map[pgtype.UUID][]RelatedCollaborator)
	idTovault := make(map[pgtype.UUID]VaultUI)
	orderedIds := []pgtype.UUID{}
	idSeen := make(map[pgtype.UUID]bool)
	var pathToPfp string
	var username string

	for _, vault := range vaultsWithCollaborators {
		vaultIdToCollaborators[vault.ID] = append(vaultIdToCollaborators[vault.ID],
			RelatedCollaborator{
				Name:      vault.CollaboratorUsername,
				Role:      string(vault.CollaboratorRole),
				PathToPfp: vault.CollaboratorPathToPfp,
			})

		if !idSeen[vault.ID] {
			orderedIds = append(orderedIds, vault.ID)
		}

		if vault.CollaboratorRole == "owner" {
			fmt.Println("VAULT NAME:", vault.Name)
			fmt.Println(vault.CollaboratorRole)
			fmt.Println(vault.CollaboratorPathToPfp)
			fmt.Println(vault.CollaboratorUsername)

			pathToPfp = vault.CollaboratorPathToPfp
			username = vault.CollaboratorUsername
		}
		idTovault[vault.ID] = NewUIvault(
			db.Vault{
				ID:          vault.ID,
				Name:        vault.Name,
				Description: vault.Description,
				Locked:      vault.Locked,
				CreatedAt:   vault.CreatedAt,
				Kind:        vault.Kind,
				UpdatedAt:   vault.UpdatedAt,
			},
			auth.Author{
				PathToPfp: pathToPfp,
				Username:  username,
			},
		)

		idSeen[vault.ID] = true
	}

	for _, id := range orderedIds {
		vault := idTovault[id]
		collaboratorsOfVault := vaultIdToCollaborators[id]

		vaultWithCollaborators := newVaultWithCollaborators(vault, collaboratorsOfVault)
		vaultsWithCollaboratorsUI = append(vaultsWithCollaboratorsUI, vaultWithCollaborators)
	}

	return vaultsWithCollaboratorsUI, nil
}

// // ---- CREATE ----
func CreateVault(queries *db.Queries, name, description string,
	userId pgtype.UUID, ctx context.Context) (*VaultWithCollaborators, error) {
	id, errCreate := queries.CreateVault(ctx, db.CreateVaultParams{
		Name:        name,
		Description: description,
		UserID:      userId,
	})
	if errCreate != nil {
		if errors.Is(errCreate, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrUnprocessable, errCreate)
	}

	vaultWithCollaborators, errGetvault := GetVaultWithCollaboratorsById(queries, userId, id, ctx)
	if errGetvault != nil {
		return nil, errGetvault
	}

	return vaultWithCollaborators, nil
}

// // ---- UPDATE ----

func UpdateVault(queries *db.Queries, vaultId, userId pgtype.UUID,
	name string, description string, locked bool, ctx context.Context) (*VaultWithCollaborators, error) {
	namePg := utils.StringToPgText(name)
	descriptionPg := utils.StringToPgText(description)
	if !namePg.Valid && !descriptionPg.Valid {
		return nil, utils.ErrNoUpdateFields
	}

	errCompletion := queries.UpdateVault(ctx, db.UpdateVaultParams{
		ID:          vaultId,
		Name:        namePg,
		Description: descriptionPg,
		Locked:      locked,
		UserID:      userId,
	})
	if errCompletion != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrUnprocessable, errCompletion)
	}

	fmt.Println(vaultId.String())
	fmt.Println(userId.String())

	vaultWithCollaborators, errGetvault := GetVaultWithCollaboratorsById(queries, userId, vaultId, ctx)
	if errGetvault != nil {
		return nil, fmt.Errorf("%w: %v", utils.ErrInternalServer, errGetvault)
	}

	return vaultWithCollaborators, nil
}

func AddCollaboratorToVault(queries *db.Queries, collaboratorUsername string,
	userId, vaultId pgtype.UUID, role string, ctx context.Context) (*VaultWithCollaborators, error) {
	err := queries.CreateCollaboratorVaultRelation(ctx, db.CreateCollaboratorVaultRelationParams{
		VaultID:              vaultId,
		CollaboratorUsername: collaboratorUsername,
		Role:                 db.VaultRelRoleType(role),
		UserID:               userId,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	vaultWithCollaborators, errvaults := GetVaultWithCollaboratorsById(queries, userId, vaultId, ctx)
	if errvaults != nil {
		return nil, errvaults
	}

	return vaultWithCollaborators, nil
}

// // ---- DELETE ----
func DeleteVault(queries *db.Queries, vaultId, userId pgtype.UUID, ctx context.Context) error {
	err := queries.DeleteVault(ctx, db.DeleteVaultParams{
		ID:     vaultId,
		UserID: userId,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
	}

	return nil
}

func RemoveCollaboratorFromVault(queries *db.Queries,
	userId, vaultId pgtype.UUID,
	collaboratorUsername string,
	ctx context.Context) (*VaultWithCollaborators, error) {
	errRelations := queries.DeleteCollaboratorVaultRelation(ctx, db.DeleteCollaboratorVaultRelationParams{
		VaultID:              vaultId,
		UserID:               userId,
		CollaboratorUsername: collaboratorUsername,
	})

	if errRelations != nil {
		if errors.Is(errRelations, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", utils.ErrBadRequest, errRelations)
	}

	vaultWithCollaborators, err := GetVaultWithCollaboratorsById(queries, userId, vaultId, ctx)
	if err != nil {
		return nil, utils.ErrBadRequest
	}

	return vaultWithCollaborators, nil
}

// // this is a one off for the vault page
// func getvaultPage(queries *db.Queries, userId, vaultId, id pgtype.UUID, ctx context.Context) (*VaultWithCollaborators, pgtype.Timestamp, error) {
// 	vaultWithRelations, err := queries.GetvaultWithTagRelations(ctx, db.GetvaultWithTagRelationsParams{
// 		ID:      id,
// 		UserID:  userId,
// 		VaultID: vaultId,
// 	})
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, pgtype.Timestamp{}, utils.ErrNotFound
// 		}
// 		return nil, pgtype.Timestamp{}, fmt.Errorf("%w: %v", utils.ErrBadRequest, err)
// 	}
// 	if len(vaultWithRelations) == 0 {
// 		return nil, pgtype.Timestamp{}, utils.ErrNotFound
// 	}
// 	var vault vaultUI
// 	collaboratorsOfVault := []RelatedCollaborator{}

// 	allCollaborators, errCollaborators := queries.GetAllCollaboratorsDesc(ctx, db.GetAllCollaboratorsDescParams{
// 		VaultID: vaultId,
// 		UserID:  userId,
// 	})
// 	if errCollaborators != nil {
// 		return nil, pgtype.Timestamp{}, fmt.Errorf("%w: %v", utils.ErrBadRequest, errCollaborators)
// 	}

// 	var createdAt pgtype.Timestamp
// 	for _, vaultWithRelation := range vaultWithRelations {
// 		vault = vaultUI{
// 			ID:   vaultWithRelation.ID.String(),
// 			Name: vaultWithRelation.Name,
// 			Idea: vaultWithRelation.Idea,
// 			Author: auth.Author{
// 				PathToPfp: vaultWithRelation.AuthorPathToPfp,
// 				Username:  vaultWithRelation.AuthorUsername,
// 			},
// 			Completed: vaultWithRelation.Completed,
// 		}
// 		createdAt = vaultWithRelation.CreatedAt

// 		if vaultWithRelation.TagID.Valid && vaultWithRelation.TagName.Valid {
// 			collaboratorsOfVault = append(collaboratorsOfVault, RelatedCollaborator{ID: vaultWithRelation.TagID.String(), Name: vaultWithRelation.TagName.String})
// 		}
// 	}

// 	availableCollaborators := getvaultAvailableCollaborators(allCollaborators, collaboratorsOfVault)
// 	vaultWithCollaborators := newVaultWithCollaborators(vault, collaboratorsOfVault, availableCollaborators)
// 	return &vaultWithCollaborators, createdAt, nil
// }
