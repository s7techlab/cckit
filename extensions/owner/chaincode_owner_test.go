package owner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe(`Chaincode owner`, func() {

	var (
		ownerIdentity1   = testdata.Certificates[0].MustIdentity(`SOME_MSP`)
		ownerIdentity2   = testdata.Certificates[1].MustIdentity(`SOME_MSP`) // same msp
		nonOwnerIdentity = testdata.Certificates[2].MustIdentity(`SOME_MSP`)

		ownerSvc = owner.NewService()
		cc, ctx  = testcc.NewTxHandler(`Chaincode owner`)
	)

	Context(`Register`, func() {

		It("Allow to register self as owner on empty state", func() {
			cc.From(ownerIdentity1).Tx(func() {
				owner, err := ownerSvc.OwnerRegister(ctx, &owner.OwnerRegisterRequest{
					MspId: ownerIdentity1.GetMSPIdentifier(),
					Cert:  ownerIdentity1.GetPEM(),
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(owner.MspId).To(Equal(ownerIdentity1.GetMSPIdentifier()))
				Expect(owner.Cert).To(Equal(ownerIdentity1.GetPEM()))
				Expect(owner.Subject).To(Equal(ownerIdentity1.GetSubject()))
				Expect(owner.Issuer).To(Equal(ownerIdentity1.GetIssuer()))
				Expect(owner.ExpiresAt.String()).To(Equal(timestamppb.New(ownerIdentity1.ExpiresAt()).String()))
				Expect(owner.UpdatedByMspId).To(Equal(ownerIdentity1.GetMSPIdentifier()))
				Expect(owner.UpdatedByCert).To(Equal(ownerIdentity1.GetPEM()))
				Expect(owner.UpdatedAt.String()).To(Equal(cc.TxTimestamp().String()))
			})
		})

		It("Disallow to register same owner once more", func() {
			cc.From(ownerIdentity1).Tx(func() {
				cc.Expect(ownerSvc.OwnerRegister(ctx, &owner.OwnerRegisterRequest{
					MspId: ownerIdentity1.GetMSPIdentifier(),
					Cert:  ownerIdentity1.GetPEM(),
				})).HasError(`state key already exists: ChaincodeOwner`)
			})
		})

		It("Disallow to register new owner from non registered identity", func() {
			cc.From(ownerIdentity2).Tx(func() {
				cc.Expect(ownerSvc.OwnerRegister(ctx, &owner.OwnerRegisterRequest{
					MspId: ownerIdentity2.GetMSPIdentifier(),
					Cert:  ownerIdentity2.GetPEM(),
				})).HasError(`tx invoker is not owner`)
			})
		})

		It("Allow to register new owner by current owner", func() {
			cc.From(ownerIdentity1).Tx(func() {
				owner, err := ownerSvc.OwnerRegister(ctx, &owner.OwnerRegisterRequest{
					MspId: ownerIdentity2.GetMSPIdentifier(),
					Cert:  ownerIdentity2.GetPEM(),
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(owner.MspId).To(Equal(ownerIdentity2.GetMSPIdentifier()))
				Expect(owner.Cert).To(Equal(ownerIdentity2.GetPEM()))
			})
		})

		It("Allow to get owner list", func() {
			cc.From(nonOwnerIdentity).Tx(func() {
				owners, err := ownerSvc.OwnersList(ctx, &emptypb.Empty{})

				Expect(err).NotTo(HaveOccurred())
				Expect(owners.Items).To(HaveLen(2))
			})
		})
	})

	Context(`Check`, func() {

		It("Non owner receives error", func() {
			cc.From(nonOwnerIdentity).Tx(func() {
				cc.Expect(ownerSvc.TxCreatorIsOwner(ctx, &emptypb.Empty{})).
					HasError(`find owner by tx creator's msp_id and cert subject: state entry not found: ChaincodeOwner`)
			})
		})

		It("Owner receives owner info", func() {
			cc.From(ownerIdentity1).Tx(func() {
				owner, err := ownerSvc.TxCreatorIsOwner(ctx, &emptypb.Empty{})

				Expect(err).NotTo(HaveOccurred())
				Expect(owner.MspId).To(Equal(ownerIdentity1.GetMSPIdentifier()))
				Expect(owner.Cert).To(Equal(ownerIdentity1.GetPEM()))
			})
		})
	})

	Context(`Update`, func() {

		It("Disallow non owner to update owner", func() {
			cc.From(nonOwnerIdentity).Tx(func() {
				cc.Expect(ownerSvc.OwnerUpdate(ctx, &owner.OwnerUpdateRequest{
					MspId: ownerIdentity2.GetMSPIdentifier(),
					Cert:  ownerIdentity2.GetPEM(),
				})).HasError(`tx invoker is not owner`)
			})
		})

		It("Disallow to update owner with same cert", func() {
			cc.From(ownerIdentity2).Tx(func() {
				cc.Expect(ownerSvc.OwnerUpdate(ctx, &owner.OwnerUpdateRequest{
					MspId: ownerIdentity2.GetMSPIdentifier(),
					Cert:  ownerIdentity2.GetPEM(),
				})).HasError(`new cert same as old cert`)
			})
		})
	})

	Context(`Delete`, func() {

		It("Disallow non owner to delete owner", func() {
			cc.From(nonOwnerIdentity).Tx(func() {
				cc.Expect(ownerSvc.OwnerDelete(ctx, &owner.OwnerId{
					MspId:   ownerIdentity2.GetMSPIdentifier(),
					Subject: ownerIdentity2.GetSubject(),
				})).HasError(`tx invoker is not owner`)
			})
		})

		It("Allow owner to delete owner, including self", func() {
			cc.From(ownerIdentity1).Tx(func() {
				deletedOwner, err := ownerSvc.OwnerDelete(ctx, &owner.OwnerId{
					MspId:   ownerIdentity1.GetMSPIdentifier(),
					Subject: ownerIdentity1.GetSubject(),
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(deletedOwner.MspId).To(Equal(ownerIdentity1.GetMSPIdentifier()))
				Expect(deletedOwner.Cert).To(Equal(ownerIdentity1.GetPEM()))
			})

			cc.From(nonOwnerIdentity).Tx(func() {
				owners, err := ownerSvc.OwnersList(ctx, &emptypb.Empty{})

				Expect(err).NotTo(HaveOccurred())
				Expect(owners.Items).To(HaveLen(1))
			})
		})

		It("Disallow to delete last owner", func() {

			cc.From(ownerIdentity2).Tx(func() {
				cc.Expect(ownerSvc.OwnerDelete(ctx, &owner.OwnerId{
					MspId:   ownerIdentity2.GetMSPIdentifier(),
					Subject: ownerIdentity2.GetSubject(),
				})).HasError(`delete last owner is not allowed`)
			})

		})
	})

})
