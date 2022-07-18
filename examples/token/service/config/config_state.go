package config

import (
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/s7techlab/cckit/router"
)

type StateService struct {
}

func NewStateService() *StateService {
	return &StateService{}
}

func (s *StateService) GetConfig(ctx router.Context, empty *emptypb.Empty) (*Config, error) {
	config, err := State(ctx).Get(&Config{}, &Config{})
	if err != nil {
		return nil, err
	}
	return config.(*Config), nil
}

// GetToken naming for token is[{TokenType}, {GroupIdPart1}, {GroupIdPart2}]
func (s *StateService) GetToken(ctx router.Context, id *TokenId) (*Token, error) {
	var (
		tokenTypeName  string
		tokenGroupName []string
		tokenType      *TokenType
		tokenGroup     *TokenGroup
		err            error
	)
	if len(id.Token) > 0 {
		tokenTypeName = id.Token[0]
	}

	if len(id.Token) > 1 {
		tokenGroupName = id.Token[1:]
	}

	tokenType, err = s.GetTokenType(ctx, &TokenTypeId{Name: tokenTypeName})
	if err != nil {
		return nil, fmt.Errorf(`token type: %w`, err)
	}

	if len(tokenGroupName) > 0 {
		tokenGroup, err = s.GetTokenGroup(ctx, &TokenGroupId{Name: tokenGroupName})
		if err != nil {
			return nil, fmt.Errorf(`token type: %w`, err)
		}
	}

	return &Token{
		Token: id.Token,
		Type:  tokenType,
		Group: tokenGroup,
	}, nil
}

func (s *StateService) CreateTokenType(ctx router.Context, req *CreateTokenTypeRequest) (*TokenType, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	tokenType := &TokenType{
		Name:        req.Name,
		Symbol:      req.Symbol,
		Decimals:    req.Decimals,
		TotalSupply: req.TotalSupply,
		GroupType:   req.GroupType,
	}

	for _, m := range req.Meta {
		tokenType.Meta = append(tokenType.Meta, &TokenMeta{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	if err := State(ctx).Insert(tokenType); err != nil {
		return nil, err
	}

	if err := Event(ctx).Set(&TokenTypeCreated{
		Name:   req.Name,
		Symbol: req.Symbol,
	}); err != nil {
		return nil, err
	}
	return tokenType, nil
}

func (s *StateService) GetTokenType(ctx router.Context, id *TokenTypeId) (*TokenType, error) {
	tokenType, err := State(ctx).Get(id, &TokenType{})
	if err != nil {
		return nil, err
	}
	return tokenType.(*TokenType), nil
}

func (s *StateService) ListTokenTypes(ctx router.Context, _ *emptypb.Empty) (*TokenTypes, error) {
	tokenTypes, err := State(ctx).List(&TokenType{})
	if err != nil {
		return nil, err
	}
	return tokenTypes.(*TokenTypes), nil
}

func (s *StateService) UpdateTokenType(ctx router.Context, request *UpdateTokenTypeRequest) (*TokenType, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StateService) DeleteTokenType(ctx router.Context, id *TokenTypeId) (*TokenType, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StateService) GetTokenGroups(ctx router.Context, id *TokenTypeId) (*TokenGroups, error) {
	tokenGroups, err := State(ctx).ListWith(&TokenGroup{}, []string{id.Name})
	if err != nil {
		return nil, err
	}
	return tokenGroups.(*TokenGroups), nil
}

func (s *StateService) CreateTokenGroup(ctx router.Context, req *CreateTokenGroupRequest) (*TokenGroup, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	_, err := s.GetTokenType(ctx, &TokenTypeId{Name: req.TokenType})
	if err != nil {
		return nil, err
	}

	tokenGroup := &TokenGroup{
		Name:        req.Name,
		TokenType:   req.TokenType,
		TotalSupply: 0,
	}

	for _, m := range req.Meta {
		tokenGroup.Meta = append(tokenGroup.Meta, &TokenMeta{
			Key:   m.Key,
			Value: m.Value,
		})
	}
	if err := State(ctx).Insert(tokenGroup); err != nil {
		return nil, err
	}

	if err := Event(ctx).Set(&TokenGroupCreated{
		Name:      req.Name,
		TokenType: req.TokenType,
	}); err != nil {
		return nil, err
	}
	return tokenGroup, nil
}

func (s *StateService) GetTokenGroup(ctx router.Context, id *TokenGroupId) (*TokenGroup, error) {
	tokenGroup, err := State(ctx).Get(id, &TokenGroup{})
	if err != nil {
		return nil, err
	}
	return tokenGroup.(*TokenGroup), nil
}

func (s *StateService) DeleteTokenGroup(ctx router.Context, id *TokenGroupId) (*Token, error) {
	//TODO implement me
	panic("implement me")
}
