package sproduct

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/RonIT-401/catalog-service/internal/app/entity"
	"github.com/RonIT-401/catalog-service/internal/app/repository/mocks"
	"github.com/RonIT-401/catalog-service/internal/pkg/testutil"
)

type createProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *createProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestCreateProductSuite(t *testing.T) {
	suite.Run(t, new(createProductSuite))
}

func (s *createProductSuite) TestNewService() {
	srvInstance := NewService(s.productRepo, s.categoryRepo)
	s.NotNil(srvInstance)
}

func (s *createProductSuite) TestCreate() {
	type args struct {
		req entity.RequestProductCreate
	}

	type want struct {
		err error
	}

	categoryGUID := uuid.Must(uuid.NewV4())

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Test Product",
					Description:  testutil.PtrString("Test Description"),
					Price:        1000,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil),
					).Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx,
						[]uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{
						{
							GUID: categoryGUID,
						},
					}, nil).
					Once()

				s.productRepo.EXPECT().
					Create(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.Name == args.req.Name &&
							p.Description == args.req.Description &&
							p.Price == args.req.Price &&
							p.CategoryGUID == args.req.CategoryGUID
					})).Return(nil).
					Once()
			},
		},
		{
			name: "already exists",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "Existing product",
					Price:        500,
					CategoryGUID: categoryGUID,
				},
			},

			want: want{
				err: entity.ErrAlreadyExists,
			},

			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil),
					).
					Return([]entity.Product{
						{Name: args.req.Name},
					}, nil).
					Once()
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductCreate{
					Name:         "New product",
					Price:        100,
					CategoryGUID: categoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					List(s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil),
					).Return([]entity.Product{}, nil).
					Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{}, nil).
					Once()
			},
		},
		{
			name: "db error on list",
			args: args{req: entity.RequestProductCreate{Name: "Err", CategoryGUID: categoryGUID}},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).Return(nil, errors.New("db error")).Once()
			},
		},
		{
			name: "db error on get categories",
			args: args{req: entity.RequestProductCreate{Name: "Err", CategoryGUID: categoryGUID}},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).Return([]entity.Product{}, nil).Once()
				s.categoryRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).Return(nil, errors.New("db error")).Once()
			},
		},
		{
			name: "db error on create",
			args: args{req: entity.RequestProductCreate{Name: "Err", CategoryGUID: categoryGUID}},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).Return([]entity.Product{}, nil).Once()
				s.categoryRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).Return([]entity.Category{{GUID: categoryGUID}}, nil).Once()
				s.productRepo.EXPECT().Create(s.ctx, mock.Anything).Return(errors.New("db error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Create(s.ctx, tc.args.req)

			if tc.want.err != nil {
				s.EqualError(err, tc.want.err.Error())
				s.Empty(result.GUID)
			} else {
				s.NoError(err)
				s.NotEmpty(result.GUID)
				s.Equal(tc.args.req.Name, result.Name)
				s.Equal(tc.args.req.Description, result.Description)
				s.Equal(tc.args.req.Price, result.Price)
				s.Equal(tc.args.req.CategoryGUID, result.CategoryGUID)
			}
		})
	}
}

type getByGUIDsProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *getByGUIDsProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestGetByGUIDsProductSuite(t *testing.T) {
	suite.Run(t, new(getByGUIDsProductSuite))
}

func (s *getByGUIDsProductSuite) TestGetByGUIDs() {

	type args struct {
		guids []uuid.UUID
	}

	type want struct {
		products []entity.Product
		err      error
	}

	productGUID := uuid.Must(uuid.NewV4())

	product := entity.Product{
		GUID: productGUID,
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "single product found",
			args: args{
				guids: []uuid.UUID{productGUID},
			},
			want: want{
				products: []entity.Product{
					{
						GUID: productGUID,
					},
				},
				err: nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, args.guids).
					Return([]entity.Product{product}, nil).Once()
			},
		},
		{
			name: "not found returns empty slice",
			args: args{
				guids: []uuid.UUID{productGUID},
			},
			want: want{
				products: []entity.Product{},
				err:      nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, args.guids).
					Return([]entity.Product{}, nil).Once()
			},
		},
		{
			name: "db error",
			args: args{guids: []uuid.UUID{productGUID}},
			want: want{products: nil, err: errors.New("db error")},
			prepare: func(args args) {
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, args.guids).Return(nil, errors.New("db error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)
			result, err := s.srv.GetByGUIDs(s.ctx, tc.args.guids)

			if tc.want.err != nil {
				s.EqualError(err, tc.want.err.Error())
				s.Empty(result)
			} else {
				s.NoError(err)
				s.Equal(tc.want.products, result)
			}
		})
	}
}

type deleteProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *deleteProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestDeleteProductSuite(t *testing.T) {
	suite.Run(t, new(deleteProductSuite))
}

func (s *deleteProductSuite) TestDelete() {

	type args struct {
		guid uuid.UUID
	}

	type want struct {
		err error
	}

	productGUID := uuid.Must(uuid.NewV4())

	product := entity.Product{
		GUID: productGUID,
	}

	deleteErr := errors.New("delete failed")

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				guid: productGUID,
			},
			want: want{
				err: nil,
			},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{product}, nil).Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, args.guid).
					Return(nil).Once()
			},
		},
		{
			name: "not found",
			args: args{
				guid: productGUID,
			},
			want: want{
				err: entity.ErrNotFound,
			},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{}, nil).Once()
			},
		},
		{
			name: "delete error",
			args: args{
				guid: productGUID,
			},
			want: want{
				err: deleteErr,
			},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).
					Return([]entity.Product{product}, nil).Once()

				s.productRepo.EXPECT().
					Delete(s.ctx, args.guid).
					Return(deleteErr).Once()
			},
		},
		{
			name: "db error on get",
			args: args{guid: productGUID},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{args.guid}).Return(nil, errors.New("db error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			err := s.srv.Delete(s.ctx, tc.args.guid)

			if tc.want.err != nil {
				s.EqualError(err, tc.want.err.Error())
			} else {
				s.NoError(err)
			}
		})
	}
}

type listProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *listProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestListProductSuite(t *testing.T) {
	suite.Run(t, new(listProductSuite))
}

func (s *listProductSuite) TestList() {
	type args struct {
		req entity.RequestProductList
	}

	type want struct {
		products []entity.Product
		err      error
	}

	categoryGUID := uuid.Must(uuid.NewV4())
	minPrice := int64(100)
	maxPrice := int64(1000)

	product := entity.Product{
		Name:  "Test Product",
		Price: 500,
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "success",
			args: args{
				req: entity.RequestProductList{
					CategoryGUID: &categoryGUID,
					MinPrice:     &minPrice,
					MaxPrice:     &maxPrice,
				},
			},
			want: want{
				products: []entity.Product{product},
				err:      nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx,
						mock.Anything,
						args.req.CategoryGUID,
						args.req.MinPrice,
						args.req.MaxPrice).
					Return([]entity.Product{product}, nil).Once()
			},
		}, {
			name: "empty result",
			args: args{
				req: entity.RequestProductList{
					CategoryGUID: &categoryGUID,
					MinPrice:     &minPrice,
					MaxPrice:     &maxPrice,
				},
			},
			want: want{
				products: []entity.Product{},
				err:      nil,
			},
			prepare: func(args args) {
				s.productRepo.EXPECT().
					List(s.ctx,
						mock.Anything,
						args.req.CategoryGUID,
						args.req.MinPrice,
						args.req.MaxPrice).
					Return([]entity.Product{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.List(s.ctx, tc.args.req)

			s.NoError(err)
			s.Equal(tc.want.products, result)
		})
	}
}

type updateProductSuite struct {
	suite.Suite
	srv          *srv
	productRepo  *mocks.MockProduct
	categoryRepo *mocks.MockCategory
	ctx          context.Context
}

func (s *updateProductSuite) SetupTest() {
	s.ctx = context.Background()
	s.productRepo = mocks.NewMockProduct(s.T())
	s.categoryRepo = mocks.NewMockCategory(s.T())
	s.srv = &srv{
		repoProduct:  s.productRepo,
		repoCategory: s.categoryRepo,
	}
}

func TestUpdateProductSuite(t *testing.T) {
	suite.Run(t, new(updateProductSuite))
}

func (s *updateProductSuite) TestUpdate() {
	type args struct {
		req entity.RequestProductUpdate
	}

	type want struct {
		err error
	}

	productGUID := uuid.Must(uuid.NewV4())
	categoryGUID := uuid.Must(uuid.NewV4())
	newCategoryGUID := uuid.Must(uuid.NewV4())

	product := entity.Product{
		GUID:         productGUID,
		Name:         "Old Product",
		Description:  testutil.PtrString("Old Description"),
		Price:        1000,
		CategoryGUID: categoryGUID,
	}

	testCases := []struct {
		name    string
		args    args
		want    want
		prepare func(args args)
	}{
		{
			name: "full update",
			args: args{
				req: entity.RequestProductUpdate{
					Name:         "New Product",
					Description:  testutil.PtrString("New Description"),
					Price:        2000,
					CategoryGUID: newCategoryGUID,
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).
					Return([]entity.Product{product}, nil).Once()

				s.productRepo.EXPECT().
					List(s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil)).
					Return([]entity.Product{}, nil).Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx,
						[]uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{
						{
							GUID: args.req.CategoryGUID,
						},
					}, nil).
					Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.GUID == productGUID &&
							p.Name == args.req.Name &&
							p.Description == args.req.Description &&
							p.Price == args.req.Price &&
							p.CategoryGUID == args.req.CategoryGUID
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "partial update - name only",
			args: args{
				req: entity.RequestProductUpdate{
					Name: "New Product",
				},
			},
			want: want{err: nil},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).
					Return([]entity.Product{product}, nil).Once()

				s.productRepo.EXPECT().
					List(s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil)).
					Return([]entity.Product{}, nil).Once()

				s.productRepo.EXPECT().
					Update(s.ctx, mock.MatchedBy(func(p entity.Product) bool {
						return p.GUID == productGUID &&
							p.Name == args.req.Name &&
							p.Description == product.Description &&
							p.Price == product.Price &&
							p.CategoryGUID == product.CategoryGUID
					})).
					Return(nil).
					Once()
			},
		},
		{
			name: "not found",
			args: args{
				req: entity.RequestProductUpdate{
					Name: "New Product",
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).
					Return([]entity.Product{}, nil).Once()
			},
		},
		{
			name: "duplicate name",
			args: args{
				req: entity.RequestProductUpdate{
					Name: "New Product",
				},
			},
			want: want{err: entity.ErrAlreadyExists},
			prepare: func(args args) {

				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).
					Return([]entity.Product{product}, nil).Once()

				s.productRepo.EXPECT().
					List(s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil)).
					Return([]entity.Product{{GUID: uuid.Must(uuid.NewV4())}}, nil).Once()
			},
		},
		{
			name: "category not found",
			args: args{
				req: entity.RequestProductUpdate{
					Name:         "New Product",
					Description:  testutil.PtrString("New Description"),
					Price:        2000,
					CategoryGUID: newCategoryGUID,
				},
			},
			want: want{err: entity.ErrNotFound},
			prepare: func(args args) {

				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).
					Return([]entity.Product{product}, nil).Once()

				s.productRepo.EXPECT().
					List(s.ctx,
						&args.req.Name,
						(*uuid.UUID)(nil),
						(*int64)(nil),
						(*int64)(nil)).
					Return([]entity.Product{}, nil).Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx,
						[]uuid.UUID{args.req.CategoryGUID}).
					Return([]entity.Category{}, nil).
					Once()
			},
		},
		{
			name: "db error on get product",
			args: args{req: entity.RequestProductUpdate{Name: "New"}},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).Return(nil, errors.New("db error")).Once()
			},
		},
		{
			name: "db error on list uniqueness",
			args: args{req: entity.RequestProductUpdate{Name: "New"}},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).Return([]entity.Product{product}, nil).Once()
				s.productRepo.EXPECT().List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).Return(nil, errors.New("db error")).Once()
			},
		},
		{
			name: "db error on update",
			args: args{req: entity.RequestProductUpdate{Name: "New"}},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().InsideTx(s.ctx, mock.Anything).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }).Once()
				s.productRepo.EXPECT().GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).Return([]entity.Product{product}, nil).Once()
				s.productRepo.EXPECT().List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).Return([]entity.Product{}, nil).Once()
				s.productRepo.EXPECT().Update(s.ctx, mock.Anything).Return(errors.New("db error")).Once()
			},
		},
		{
			name: "db error on get category during update",
			args: args{
				req: entity.RequestProductUpdate{
					Name:         "New Product",
					CategoryGUID: newCategoryGUID,
				},
			},
			want: want{err: errors.New("db error")},
			prepare: func(args args) {
				s.categoryRepo.EXPECT().
					InsideTx(s.ctx, mock.AnythingOfType("func(context.Context) error")).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Once()

				s.productRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{productGUID}).
					Return([]entity.Product{product}, nil).Once()

				s.productRepo.EXPECT().
					List(s.ctx, &args.req.Name, (*uuid.UUID)(nil), (*int64)(nil), (*int64)(nil)).
					Return([]entity.Product{}, nil).Once()

				s.categoryRepo.EXPECT().
					GetByGUIDs(s.ctx, []uuid.UUID{args.req.CategoryGUID}).
					Return(nil, errors.New("db error")).
					Once()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.prepare(tc.args)

			result, err := s.srv.Update(s.ctx, productGUID, tc.args.req)

			if tc.want.err != nil {
				s.EqualError(err, tc.want.err.Error())
				s.Empty(result.GUID)
			} else {
				s.NoError(err)
				s.NotEmpty(result.GUID)
				s.Equal(productGUID, result.GUID)

				if tc.args.req.Name != "" {
					s.Equal(tc.args.req.Name, result.Name)
				} else {
					s.Equal(product.Name, result.Name)
				}

				if tc.args.req.Description != nil {
					s.Equal(tc.args.req.Description, result.Description)
				} else {
					s.Equal(product.Description, result.Description)
				}

				if tc.args.req.Price != 0 {
					s.Equal(tc.args.req.Price, result.Price)
				} else {
					s.Equal(product.Price, result.Price)
				}

				if tc.args.req.CategoryGUID != uuid.Nil {
					s.Equal(tc.args.req.CategoryGUID, result.CategoryGUID)
				} else {
					s.Equal(product.CategoryGUID, result.CategoryGUID)
				}
			}
		})
	}
}
